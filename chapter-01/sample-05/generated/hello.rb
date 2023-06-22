#!/usr/bin/ruby
require "matrix"

def rand_number(max_value):
    return (Mat[rand(6).to_i, 6] * max_value)
end
puts rand_number(6) # Prints a random number between 1 and 6