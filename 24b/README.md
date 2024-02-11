This one was tricky. Tried solving it via a simple self-made CAS. However this took too long for not a clear answer as how to solve it generally.

I then used python with sympy package. Here I tried to solve it step-by-step. This yielded no answer as well.

By further analysing the sample input I stumbled upon the fact that one of the solutions makes one of the equations zero. These is the approach to test the valid starts.

A "valid start" is that one of the hailstones has the same coordinate (x, y, or z) and its velocity the same as the starting shooting position.

Then further testing if all the 900 equations still hold (using the assumptions of the variables, i.e. the positions are integers, and the times t0 to t299 are nonzero and positive reals)

This let to the only solution using this system, which was the answer to the coding question.
