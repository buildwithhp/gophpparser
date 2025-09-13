<?php
// Phase 1 test features

// While loop
$i = 0;
while ($i < 10) {
    echo "Count: " . $i;
    $i++;
}

// Foreach loop with associative array
$users = ["john" => 25, "jane" => 30, "bob" => 35];
foreach ($users as $name => $age) {
    echo "User: " . $name . ", Age: " . $age;
}

// Break and continue
for ($j = 0; $j < 20; $j++) {
    if ($j == 5) {
        continue;
    }
    if ($j == 15) {
        break;
    }
    echo $j;
}

// Regular indexed array
$numbers = [1, 2, 3, 4, 5];
foreach ($numbers as $num) {
    echo $num;
}
?>