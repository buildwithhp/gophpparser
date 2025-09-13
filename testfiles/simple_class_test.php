<?php
class User {
    public $name;
    private $age;
    
    public function getName() {
        return $this->name;
    }
}

$user = new User();
$name = $user->getName();
?>