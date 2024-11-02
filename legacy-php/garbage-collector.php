<?php
/**
 * Garbage Collector
 * Used to delete files created earlier than N days in a certain directory, including subdirectories.
 *
 * PLEASE BE CAREFUL! AUTHORS ARE NOT RESPONSIBLE IF YOU ACCIDENTALLY DELETE SOMETHING IMPORTANT!
 *
 * @copyright   Copyright (c) 2024, https://github.com/utilmind/
 * @license     MIT License (https://opensource.org/licenses/MIT)
 * @version     1.0.0
 * @since       5.4
 */

// gettings arguments
if (!is_array($argv) || count($argv) < 2) {
    echo <<<END
Usage: $argv[0] [directory] [expire after N days (optional, default 90)]
Example: $argv[0] /var/www/project/temp-dir 90

END;
    exit;
}


// Alternative is: system('rm -rf ' . escapeshellarg($dir), $retval);
function rmdirr($dir,
                $expire_time = null, // in seconds. If $expire_time specified, only files modified more than N seconds ago will be deleted
                $subdirs_only = null, // skip files in the $dir. Focus only on content of subdirectories.
                $debug_echo = null) {

    if (strlen($dir) < 6) { // AK: !!WARNING!! I have killed the half of my HDD 2012-09-07 and hardly restored with "FreeUndelete" utility!
        die('Hey, I just saved your life!'); // we don't allow deletion of the top-level directories
    }

    if (file_exists($dir = rtrim($dir, '/\\ ') . DIRECTORY_SEPARATOR)) { // add trailing /. Need to append the filename to path.
        $skip = false;
        if ($objs = scandir($dir)) {
            foreach($objs as &$obj) {
                if ($obj !== '.' && $obj !== '..') {
                    $o = $dir.$obj;

                    if (is_dir($o)) {
                        rmdirr($o, $expire_time, null, $debug_echo); // Recursion!
                    }elseif (!$subdirs_only) {
                        if ($expire_time && (@filemtime($o) > $_SERVER['REQUEST_TIME'] - $expire_time)) {
                            $skip = true;
                            continue;
                        }
                        if ($debug_echo) {
                            echo "Unlink $o\n";
                            unlink($o);
                        }else {
                            @unlink($o);
                        }
                    }
                }
            }
        }
        if (!$subdirs_only && !$skip) { // remove only if there is no files inside.
            @rmdir($dir); // AK: 2 errors are possible here: directory not empty (legal to not delete) OR permission denied. PLEASE WATCH DIRECTORY PRIVILEGES!
        }
    }

    return $dir;
}


$days = isset($argv[2]) ? (int)$argv[2] : 0;
if (0 >= $days) { // don't accept negative values, use default, like it was omitted
    $days = 90;
}
rmdirr($argv[1], $days * 60 * 60 * 24, true, true);
