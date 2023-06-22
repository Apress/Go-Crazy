package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apiResource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sort"

	learnv1alpha1 "github.com/gocrazy/blockchain-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BlockchainReconciler reconciles a Blockchain object
type BlockchainReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=learn.gocrazy.com,resources=blockchains,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=learn.gocrazy.com,resources=blockchains/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=learn.gocrazy.com,resources=blockchains/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *BlockchainReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log.SetPrefix("BlockchainReconciler")
	blockchain := &learnv1alpha1.Blockchain{}
	err := r.Get(ctx, req.NamespacedName, blockchain)

	if err != nil {
		return reconcile.Result{}, err
	}

	/*
		log.Println("namespace", blockchain.Namespace, blockchain.GetNamespace(), req.NamespacedName)
		log.Println("name", blockchain.Name)
		log.Println("replicas", *blockchain.Spec.Replicas)
		log.Println("image", blockchain.Spec.Image)
		log.Println("cpu", blockchain.Spec.Cpu)
		log.Println("memory", blockchain.Spec.Memory)
		log.Println("api", blockchain.Spec.ApiPort)

		for _, value := range blockchain.Spec.Command {
			log.Printf("command %s\n", value)
		}

		for _, value := range blockchain.Spec.ClientArgs {
			log.Printf("ClientArgs %s\n", value)
		}
	*/

	// Check if the statefulset already exists, if not create a new one
	foundSts := &appsv1.StatefulSet{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: blockchain.Name, Namespace: blockchain.Namespace}, foundSts)
	if err != nil && errors.IsNotFound(err) {
		// Create a new StatefulSet
		sts := r.ReconcileStatefulSet(blockchain)

		err = r.Client.Create(context.TODO(), sts)
		if err != nil {
			log.Println("Failed to create new StatefulSet", err, "Namespace", sts.Namespace, "Name", sts.Name)
			return reconcile.Result{}, err
		}
		// StatefulSet created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Println("Failed to get StatefulSet", err)
		return reconcile.Result{}, err
	}

	// sts already exists. Updating to reflect the Blockchain spec
	// Ensure the number of replicas matches the spec
	r.reconcileReplicas(blockchain, foundSts)
	// Ensure the container image size is the same as the spec
	r.reconcileImage(blockchain, foundSts)
	// Ensure the ClientArgs are the same as the spec. The order does not matter
	r.reconcileArgs(blockchain, foundSts)
	// Ensure the command is the same as the spec. The order does matter
	r.reconcileCommand(blockchain, foundSts)
	// Ensure the resources is the same as the spec.
	r.reconcileResources(blockchain, foundSts)
	// Ensure the container ports are the same as the spec.
	// r.reconcileContainerPorts(blockchain, foundSts)

	return ctrl.Result{}, nil
}

func (r *BlockchainReconciler) ReconcileStatefulSet(b *learnv1alpha1.Blockchain) *appsv1.StatefulSet {
	log.Println("Creating a new StatefulSet")

	// Make sure to run at least 1 replicas
	if b.Spec.Replicas == nil {
		b.Spec.Replicas = pointer.Int32(1)
	}

	// provisioning a PVC to store this statefulset's data
	pvc := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "data",
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			StorageClassName: pointer.String("standard"),
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: apiResource.MustParse(fmt.Sprintf("%dGi", 1)),
				},
			},
		},
	}

	// Specifying default resources for the main container
	if b.Spec.Cpu == "" {
		b.Spec.Cpu = "500m"
	}

	if b.Spec.Memory == "" {
		b.Spec.Cpu = "1Gi"
	}

	reqs := &v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu":    apiResource.MustParse(b.Spec.Cpu),
			"memory": apiResource.MustParse(b.Spec.Memory),
		},
		Requests: v1.ResourceList{
			"cpu":    apiResource.MustParse(b.Spec.Cpu),
			"memory": apiResource.MustParse(b.Spec.Memory),
		},
	}

	// specify default value for the api port
	if b.Spec.ApiPort == 0 {
		b.Spec.ApiPort = 8545
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.Name,
			Namespace: b.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: b.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: b.ObjectMeta.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.ObjectMeta.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           b.Spec.Image,
						ImagePullPolicy: "Always",
						Name:            "app",
						Command:         b.Spec.Command,
						Args:            b.Spec.ClientArgs,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 30303,
							Name:          "p2p",
							Protocol:      "TCP",
						}, {
							ContainerPort: b.Spec.ApiPort,
							Name:          "api",
							Protocol:      "TCP",
						}},
						Resources: *reqs,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "data",
							MountPath: "/data",
						}},
					}},
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				pvc,
			},
		},
	}
	// Set Learn instance as the owner and controller
	controllerutil.SetControllerReference(b, sts, r.Scheme)
	return sts
}

func (r *BlockchainReconciler) reconcileReplicas(blockchain *learnv1alpha1.Blockchain, sts *appsv1.StatefulSet) (ctrl.Result, error) {
	specReplicas := *blockchain.Spec.Replicas
	stsReplicas := *sts.Spec.Replicas
	log.Println("Sts found. replicas", specReplicas, "foundSts.Spec.Replicas", stsReplicas)

	if stsReplicas != specReplicas {
		sts.Spec.Replicas = &specReplicas
		err := r.Client.Update(context.TODO(), sts)
		if err != nil {
			log.Println("Failed to update StatefulSet Replicas", err, "Namespace", sts.Namespace, "Name", sts.Name)
			return reconcile.Result{}, err
		}
		log.Println("Spec updated", "Replicas", *sts.Spec.Replicas)
	}

	// Spec unchanged or updated - return and requeue
	return reconcile.Result{Requeue: true}, nil
}

func (r *BlockchainReconciler) reconcileImage(blockchain *learnv1alpha1.Blockchain, sts *appsv1.StatefulSet) (ctrl.Result, error) {
	specImage := blockchain.Spec.Image
	stsImage := sts.Spec.Template.Spec.Containers[0].Image
	if stsImage != specImage {
		sts.Spec.Template.Spec.Containers[0].Image = specImage
		err := r.Client.Update(context.TODO(), sts)
		if err != nil {
			log.Println("Failed to update StatefulSet Image", err, "Namespace", sts.Namespace, "Name", sts.Name)
			return reconcile.Result{}, err
		}
		log.Println("Spec updated", "Image", sts.Spec.Template.Spec.Containers[0].Image)
	}

	// Spec unchanged or updated - return and requeue
	return reconcile.Result{Requeue: true}, nil
}

func (r *BlockchainReconciler) reconcileArgs(blockchain *learnv1alpha1.Blockchain, sts *appsv1.StatefulSet) (ctrl.Result, error) {
	specClientArgs := blockchain.Spec.ClientArgs
	stsClientArgs := sts.Spec.Template.Spec.Containers[0].Args

	argsEquals := compareSlices(specClientArgs, stsClientArgs, true)
	log.Println("reconcileArgs", "argsEquals", argsEquals)

	if !argsEquals {
		sts.Spec.Template.Spec.Containers[0].Args = specClientArgs
		err := r.Client.Update(context.TODO(), sts)
		if err != nil {
			log.Println("Failed to update StatefulSet ClientArgs", err, "Namespace", sts.Namespace, "Name", sts.Name)
			return reconcile.Result{}, err
		}
		log.Println("Spec updated", "ClientArgs", sts.Spec.Template.Spec.Containers[0].Args)
	}
	// Spec unchanged or updated - return and requeue
	return reconcile.Result{Requeue: true}, nil
}

func (r *BlockchainReconciler) reconcileCommand(blockchain *learnv1alpha1.Blockchain, sts *appsv1.StatefulSet) (ctrl.Result, error) {
	specCommand := blockchain.Spec.Command
	stsContainerCommand := sts.Spec.Template.Spec.Containers[0].Command

	argsEquals := compareSlices(specCommand, stsContainerCommand, false)
	log.Println("reconcileCommand", "argsEquals", argsEquals)

	if !argsEquals {
		sts.Spec.Template.Spec.Containers[0].Command = specCommand
		err := r.Client.Update(context.TODO(), sts)
		if err != nil {
			log.Println("Failed to update StatefulSet Command", err, "Namespace", sts.Namespace, "Name", sts.Name)
			return reconcile.Result{}, err
		}
		log.Println("Spec updated", "Command", sts.Spec.Template.Spec.Containers[0].Command)
	}
	// Spec unchanged or updated - return and requeue
	return reconcile.Result{Requeue: true}, nil
}

func (r *BlockchainReconciler) reconcileResources(blockchain *learnv1alpha1.Blockchain, sts *appsv1.StatefulSet) (ctrl.Result, error) {

	specCpu := blockchain.Spec.Cpu
	stsContainerResources := sts.Spec.Template.Spec.Containers[0].Resources
	stsResourceRequestCpu := stsContainerResources.Requests.Cpu().String()

	specMemory := blockchain.Spec.Memory
	stsResourceRequestMemory := stsContainerResources.Requests.Memory().String()

	if specCpu != stsResourceRequestCpu || specMemory != stsResourceRequestMemory {

		reqs := &v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    apiResource.MustParse(specCpu),
				"memory": apiResource.MustParse(specMemory),
			},
			Requests: v1.ResourceList{
				"cpu":    apiResource.MustParse(specCpu),
				"memory": apiResource.MustParse(specMemory),
			},
		}

		sts.Spec.Template.Spec.Containers[0].Resources = *reqs

		err := r.Client.Update(context.TODO(), sts)
		if err != nil {
			log.Println("Failed to update StatefulSet Resources", err, "Namespace", sts.Namespace, "Name", sts.Name)
			return reconcile.Result{}, err
		}
		log.Println("Spec updated", "Resources", sts.Spec.Template.Spec.Containers[0].Resources.String())
	}
	// Spec unchanged or updated - return and requeue
	return reconcile.Result{Requeue: true}, nil
}

/*
func (r *BlockchainReconciler) reconcileContainerPorts(blockchain *learnv1alpha1.Blockchain, sts *appsv1.StatefulSet) (ctrl.Result, error) {

	specApiPort := blockchain.Spec.ApiPort

	stsReplicas := sts.Spec.Template.Spec.Containers[0].Ports[0].
	log.Println("Sts found. replicas", specReplicas, "foundSts.Spec.Replicas", stsReplicas)

	if stsReplicas != specReplicas {
		sts.Spec.Replicas = &specReplicas
		err := r.Client.Update(context.TODO(), sts)
		if err != nil {
			log.Println("Failed to update StatefulSet Replicas", err, "Namespace", sts.Namespace, "Name", sts.Name)
			return reconcile.Result{}, err
		}
		log.Println("Spec updated", "Replicas", *sts.Spec.Replicas)
	}

	// Spec unchanged or updated - return and requeue
	return reconcile.Result{Requeue: true}, nil
}
*/

// utility function to compare 2 slices for equality.
func compareSlices(s1 []string, s2 []string, withPreOrdering bool) bool {
	if len(s1) != len(s2) {
		return false
	}

	// Sort the slices so that their elements are in the same order
	if withPreOrdering {
		sort.Strings(s1)
		sort.Strings(s2)
	}

	// Compare the elements of the sorted slices
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlockchainReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&learnv1alpha1.Blockchain{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
