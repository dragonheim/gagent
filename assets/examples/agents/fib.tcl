##########################################
### Perform Fibanaci sequence up to 10 ###
##########################################
set GHINT [split "math, fib" ,]
proc fib {x} {
	if {<= $x 1} {
		return 1
	} else {
		+ [fib [- $x 1]] [fib [- $x 2]]
	}
}

puts [fib 20]