local i = 0
print("test goto")

goto a1
:: a2
::
i=i+1
goto a3

::a1::
assert(true)
i=i+1
goto a2

:: a3::
assert(i == 2)

if false then
::iflabel::
    i = 10
end

if i== 10 then goto iflabel2 end

goto iflabel ::iflabel2::

assert(i == 10)

function a()
::innerlabel::
assert(false)
end

-- goto innerlabel panic

_GOTO = "goto"
