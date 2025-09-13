<?php
function greet($name) {
    echo "Hello, " . $name . "\n";
    return true;
}

$message = "World";
$result = greet($message);

if ($result) {
    echo "Greeting successful";
} else {
    echo "Greeting failed";
}

$numbers = [1, 2, 3, 4, 5];
$sum = 0;

function sumAndGreet($x)
{
    global $sum;
    for ($i = 0; $i < $x; $i++) {
        $sum = $sum + $i;
        greet("User " . $sum);
    }
}

sumAndGreet(10);

echo "Sum: " . $sum;
?>