<?php
$i = 0;
while ($i < 10) {
    echo $i;
    $i++;
}

$users = ["john" => 25, "jane" => 30];
foreach ($users as $name => $age) {
    echo $name;
}

for ($j = 0; $j < 5; $j++) {
    if ($j == 2) {
        continue;
    }
    if ($j == 4) {
        break;
    }
    echo $j;
}
?>