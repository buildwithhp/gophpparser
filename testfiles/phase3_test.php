<?php
namespace App;

use Models\User;

try {
    $user = new User();
    $result = $user->process();
} catch (Exception $e) {
    throw new CustomException("Failed: " . $e->getMessage());
} finally {
    cleanup();
}

$callback = function($x, $y) use ($multiplier) {
    yield $x * $y * $multiplier;
};
?>