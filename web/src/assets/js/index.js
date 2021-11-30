var PERSON_IMG = "../assets/img/iu2.jpeg";
var PERSON_NAME = "아이유(IU) 공식 트위터";
var PERSON_TAG = "@_IUofficial";

var feed = document.getElementsByClassName("posts")
var text = document.getElementById("msg");
var send = document.getElementById("send")

send.onclick = function (e) {
  if (text.value != ""){
    handleMessageEvent()
  }
}

function handleMessageEvent() {
  msg = getMessage(PERSON_NAME, PERSON_IMG, PERSON_TAG, text.value, "");
  insertMsg(msg, feed[0]);
  text.value = "";
}


function getMessage(profilename, profileimg, username, text, img) {
  const d = new Date()
  var msg = `
  <div class="post">
    <div class="post__avatar">
      <img src="${profileimg}" alt=""/>
    </div>
    <div class="post__body">
      <div class="post__header">
        <div class="post__headerText">
          <h3>${profilename} <span class="post__headerSpecial"> <span class="material-icons post__badge"> verified </span>${username}</span></h3>
        </div>
        <div class="post__headerDescription">
          <span style="white-space: pre-line">${text}</span>
        </div>
      </div>
  `
  if (img != "") {
    msg = msg + `<img src="${img}" alt=""/>`
  }
  msg = msg + `</div></div>`
  return msg;
}

function insertMsg(msg, domObj) {
  domObj.insertAdjacentHTML("afterbegin", msg);
  domObj.scrollTop += 500;
}