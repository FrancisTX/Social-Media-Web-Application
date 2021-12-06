<?php
if (isset($_POST['submit'])) {
    //$username = $_POST['username'];
    //$followusername = $_POST['followusername'];
    
    $domdoc = new DOMDocument();
    $domdoc->loadHTML("../html/search.html");
    
    $followusername = $domdoc->getElementById('follow')->nodeValue;

    $host = "127.0.0.1:3306";
    $dbUsername = "root";
    $dbPassword = "wtx20150914";
    $dbName = "twitter";

    $conn = new mysqli($host, $dbUsername, $dbPassword, $dbName);

    if ($conn->connect_error) {
        console.log('Could not connect to the database.');
    }
    else {
        $Select = "SELECT * FROM followers(username, followusername) values(?, ?)";

        $stmt = $conn->prepare($Select);
        $stmt->bind_param("ss", $username, $followusername);
        $stmt->execute();
        $stmt->bind_result($resultFollow);
        $stmt->store_result();
        $stmt->fetch();
        $rnum = $stmt->num_rows;
        if ($rnum == 0) {
            $stmt->close();
            $Insert = "INSERT INTO followers(username, followusername) values(?, ?)";

            $stmt = $conn->prepare($Insert);
            $stmt->bind_param("ss",$username, $followusername);
            if ($stmt->execute()) {
                console.log("Follow sucessfully.");
            }
            else {
                console.log($stmt->error);
            }
        }
        else{
            console.log("Already Followed");
        }
    }
    $stmt->close();
    $conn->close();
}
else {
    console.log("Submit button is not set");
}
?>