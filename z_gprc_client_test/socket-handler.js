let socket;

// function login(username, phone){
//     const xhttp = new XMLHttpRequest();
//     xhttp.onload = function() {
//         console.log(this.responseText);
//         if (this.status === 200){
//             const result = JSON.parse(this.responseText);
//             if (result.status === 200){
//                 initSocket(username, phone);
//             }else{
//                 console.log(result.message);
//             }
//
//         }
//
//     }
//
//     let body = `{"username": "${username}","password": "${password}"}`;
//     console.log("body",body);
//
//     xhttp.open("POST", "https://api.cwm.com/subacc/loginAccSub", true);
//     // xhttp.open("POST", "http://localhost:6868/subacc/loginAccSub", true);
//     xhttp.setRequestHeader("Content-type", "application/json");
//     xhttp.send(body);
// }

function login(username, phone){
    initSocket(username, phone);
}

function initSocket(username, phone) {
    // socket = io("https://api.cwm.com", {
    //     path: '/ws',
    //     transports: ['websocket'],
    //     secure: true,
    //     query: {
    //         token: token
    //     }
    // });

    socket = io("http://localhost:9000", {
        path: '/ws',
        transports: ['websocket'],
        query: {
            username: username,
            phone: phone,
        }
    });


    socket.on("connect_error", (err) => {
        console.log(err.message); // prints the message associated with the error
        socket.disconnect();
    });

    socket.on("hi", (msg) => {
        console.log(msg); // prints the message associated with the error
    });

    socket.on("connect", () => {
        console.log("socket connected");
        form.style.visibility = "hidden";
        form2.style.visibility = "visible";

    });

    socket.on("disconnect", () => {
        console.log("socket disconnected");
    });

    socket.on('onChatMsg', function(msg) {
        let item = document.createElement('li');
        // item.textContent = JSON.stringify(msg);
        item.textContent = msg;
        messages.appendChild(item);
        window.scrollTo(0, document.body.scrollHeight);
    });

    socket.on('onInteruptSession', function(msg) {
        console.log('onInteruptSession:',msg);
    });
}


function sendMsg(msg){

    socket.emit("test", msg, (response) => {
        msgContent.value = "";
        console.log(response);
        // if (response.status === 200){
        //     console.log("Success!");
        // }else{
        //     console.log("Error:",response.status, response.message);
        // }
    });

    // socket.emit("chat", msg, (response) => {
    //     msgContent.value = "";
    //     console.log(response);
    //     // if (response.status === 200){
    //     //     console.log("Success!");
    //     // }else{
    //     //     console.log("Error:",response.status, response.message);
    //     // }
    // });

    // socket.emit("chat2", {msg: msg}, (response) => {
    //     msgContent.value = "";
    //     console.log(response);
    //     if (response.status === 200){
    //         console.log("Success!");
    //     }else{
    //         console.log("Error:",response.status, response.message);
    //     }
    // });
}