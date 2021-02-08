/*
    https://stackoverflow.com/questions/63780807/calculate-nth-fibonacci-number-using-risc-v-rv32i-compiler-without-recursion
*/

li a0, 20

fib:
    li a1, 0             # This is a
    li a2, 1             # This is b
    li a3, 0             # This is c 
    li a4, 1             # This is i
   
    li a6, 2             # dummy, just to check a0 with
    ble a0, a6, cond1    # check if a0 <= 2
    bgt a0, a6, head     # if a0 > 2, proceed to loop

head:                    # start of loop
    add a3, a1, a2       # Here I'm implementing the space optimized version of fibonacci series without recursion:
    mv a1, a2            # for i in range(2, n+1): c = a + b; a = b; b = c; return b
    mv a2, a3            
    addi a4, a4, 1
    blt a4, a0, head
    bge a4, a0, end      # iterates n-1 times and goes to end

cond1:
    li a0, 1             # if n==1 or n==2 return 1
    li a7, 93
    
end:
    mv a0, a2            # copying the value of a2 (which is b) to a0, since the testbench
    li a7, 93            # runtest.s is setup that way.
    