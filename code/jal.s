main:
  li x29, 1
  jal x4, target
  li x29, 2
  li x29, 3
  li x29, 4

target:
  li x29, 10
  jal x4, target
