<?php
// Ternary operator
$result = $condition ? "true" : "false";

// Spaceship operator  
$comparison = $a <=> $b;

// Null coalescing
$username = $input ?? "default";

// Null coalescing assignment
$config ??= [];

// Nullsafe operator
$length = $user?->getName()?->length();
?>