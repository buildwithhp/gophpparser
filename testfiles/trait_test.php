<?php
trait Loggable {
    public function log($message) {
        echo "Log: " . $message;
    }
}

class User implements UserInterface {
    use Loggable;
    
    const STATUS_ACTIVE = 1;
    private $name;
    
    public function getName() {
        return $this->name;
    }
}
?>