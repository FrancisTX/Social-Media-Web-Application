<?php
if (isset($_POST['submit'])) {
    //$username = $_POST['username'];
    //$unfollowusername = $_POST['unfollowusername'];
    $domdoc = new DOMDocument();
    $domdoc->loadHTML("../html/search.html");
    
    $followusername = $domdoc->getElementById('unfollow')->nodeValue;

    $host = "127.0.0.1:3306";
    $dbUsername = "root";
    $dbPassword = "wtx20150914";
    $dbName = "twitter";

    $conn = new mysqli($host, $dbUsername, $dbPassword, $dbName);

    if ($conn->connect_error) {
        console.log('Could not connect to the database.');
    }
    else {
        $Select = "SELECT * FROM followers(username, unfollowusername) values(?, ?)";

        $stmt = $conn->prepare($Select);
        $stmt->bind_param("ss", $username, $unfollowusername);
        $stmt->execute();
        $stmt->bind_result($resultFollow);
        $stmt->store_result();
        $stmt->fetch();
        $rnum = $stmt->num_rows;
        if ($rnum == 0) {
            console.log("Not followed yet");
        }
        else{
            $stmt->close();
            $Delete = "DELETE FROM followers(username, unfollowusername) values(?, ?)";

            $stmt = $conn->prepare($Delete);
            $stmt->bind_param("ss",$username, $unfollowusername);
            if ($stmt->execute()) {
                console.log("Unfollow sucessfully.");
            }
            else {
                console.log($stmt->error);
            }
        }
    }
    $stmt->close();
    $conn->close();
}
else {
    console.log("Submit button is not set");
}
?>