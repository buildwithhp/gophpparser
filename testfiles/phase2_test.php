<?php
// Phase 2 OOP test features

class User {
    public $name;
    private $age;
    protected $email;
    
    public function __construct($name, $age) {
        $this->name = $name;
        $this->age = $age;
    }
    
    public function getName() {
        return $this->name;
    }
    
    private function getAge() {
        return $this->age;
    }
    
    public static function validate($data) {
        return true;
    }
}

class Admin extends User {
    public $permissions;
    
    public function __construct($name, $age, $permissions) {
        $this->name = $name;
        $this->age = $age;
        $this->permissions = $permissions;
    }
    
    public function hasPermission($permission) {
        return true;
    }
}

// Object instantiation
$user = new User("John", 30);
$admin = new Admin("Jane", 25, ["read", "write"]);

// Object method calls
$userName = $user->getName();
$isValid = User::validate($userData);

// Property access
$name = $user->name;
$permissions = $admin->permissions;
?>