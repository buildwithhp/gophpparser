<?php
class User {
    public static function validate($data) {
        return true;
    }
}

$result = User::validate($data);
?>