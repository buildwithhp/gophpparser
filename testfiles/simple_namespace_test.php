<?php
namespace App;

use Models;

try {
    $result = risky_operation();
} catch (Exception $e) {
    echo $e->getMessage();
}

$callback = function($x, $y) {
    return $x + $y;
};
?>