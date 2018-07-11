import React, { Component } from 'react';
import { Grid, Container,Menu,Button, Icon, Segment, Modal,Input,Dropdown} from 'semantic-ui-react'
import UserList from './UserList';
import ChatMessageBox from './ChatMessageBox'
import * as moment from 'moment';
import cookies from 'js-cookie';

class Chatroom extends Component{

  constructor(props){
    super(props)
    //var guest = "guest_"+this.makeid()
    this.state = {
      message:"",
      messageList:[],
      guestname:"guest_"+this.makeid(),
      username:"",
      userList:[],
      currentRoom:"",
      currentRoomId:"",
      newRoomName:"",
      roomOptions:[],
      modalOpen:false,
      roomModal:false,
    }
   // this.setState({messageList:[]})
    this.initSocket = this.initSocket.bind(this);
    this.sendMessage = this.sendMessage.bind(this);
    this.loginFormOpen =this.loginFormOpen.bind(this)
    this.setUsername = this.setUsername.bind(this)
    this.close = this.close.bind(this)
    this.logIn = this.logIn.bind(this)
    this.newRoomOpen = this.newRoomOpen.bind(this)
    this.setNewRoomname = this.setNewRoomname.bind(this)
    this.createNewRoom = this.createNewRoom.bind(this)
    this.setCurrentRoom = this.setCurrentRoom.bind(this)
  }

   makeid() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 6; i++)
      text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
  }

  newRoomOpen(event){
    event.preventDefault();
    this.setState({roomModal:true})
  }

  setNewRoomname(event){
    event.preventDefault()
    this.setState({newRoomName:event.target.value})
  }

  createNewRoom(){
    fetch("/api/v1/rooms/",{
      method:"POST",
      body: JSON.stringify({
        text:this.state.newRoomName,
        value:"r"+this.makeid(),
      }),
      headers:{
        'Content-Type': 'application/json'
      }
    }).then((res) => {
      return res.json();
    }).then((res)=>{
      console.log(res)
      let rs = this.state.roomOptions
      rs.push(res)
      this.setState({currentRoom:res.value,roomModal:false,currentRoomId:res.key})
      console.log(this.state.currentRoomId)
    }).catch(error => console.error('Error:', error))
  }

  setCurrentRoom(event,data){
    event.preventDefault()
    this.setState({currentRoom:data})
  }

  setUsername(event){
    event.preventDefault()
    this.setState({username:event.target.value})
  }
  initSocket () {
    var hostname=(window.location.hostname)
    this.ws = new WebSocket("ws://"+hostname+":18000/ws");
    this.ws.onmessage = (msg) => {
      var prd_msg = JSON.parse(msg.data)
      this.setState({ message: msg.data });
     if (prd_msg[0].message_type==='user_list'){
       this.setState({userList:prd_msg[0].list})
       console.log(this.state.userList)
     }else{
        if (prd_msg[0].message!==""){
            this.setState({ messageList: [...this.state.messageList, prd_msg[0]] })
            cookies.set('LastMessageId', prd_msg[0].message_id, { path: '/' });
            console.log(prd_msg[0])
        }
      }
    }
    console.log("init ws connection")
    this.ws.onopen =()=>{
      console.log("connected")
    }
  }
  componentDidMount(){
    this.initSocket()
    fetch("/api/v1/messages/?LastMessageId="+cookies.get("LastMessageId")).then((res) => {
      return res.json();
    }).then((res) => {
      if (res.length>0){
        cookies.set('LastMessageId', res[res.length-1].message_id, { path: '/' });
      }
      this.setState({messageList:res})
    }).catch((err) => {
      this.setState({err});
    });

  //  Cookies.set('test-cookies', "hello test cookies", { path: '/' });

    fetch("/api/v1/rooms/").then((res) => {
      return res.json();
    }).then((res) => {
       // console.log(this.state.messages)
      this.setState({roomOptions:res})
      //this.roomOptions = res
    //s  console.log(this.state.roomOptions)
    }).catch((err) => {
      this.setState({err});
    });
  }

  generateTimestamp () {
   // console.log(moment().format("MMMM DD YYYY, h:mm:ss a"))
    //var iso = new Date().toTimeString() //.toISOString();
    return moment().format("YYYY-MM-DD h:mm:ss");
  }

  sendMessage (message) {
    console.log(this.state.currentRoomId)
    if (this.state.currentRoom ===""){
      alert("Please select a chat room!")
      return
    }
      this.ws.send(
        JSON.stringify({
          username: this.state.username===""?this.state.guestname:this.state.username,
          register:false,
          guestname:this.state.guestname,
          message: (message+" ("+this.generateTimestamp()+")"),
          room:this.state.currentRoom,
          room_id:this.state.currentRoomId,
        })
      );
  }
  logIn(){
    this.ws.send(
      JSON.stringify({
        username:this.state.username,
        guestname:this.state.guestname,
        register:true,
        message:"",
        room:this.state.currentRoom,
        room_id:this.state.currentRoomId,
      })
    )
    this.setState({modalOpen:false})
  }

  loginFormOpen(event){
    event.preventDefault();
    this.setState({ modalOpen: true})
  }
  close(event){
    this.setState({ modalOpen: false,roomModal:false })
  }
  render(){
    return (
    <Container className="chatBoard">
    <div className="ui ">
          <Menu icon size='tiny'>
            <Menu.Item name='side layout' active >
              <Icon name='sidebar' />
            </Menu.Item>
            <Menu.Item name='side layout' active >
              Chat Room
            </Menu.Item>
            <Menu.Item name='side layout' active >
              <Button icon labelPosition='left' color="teal" onClick={this.newRoomOpen}>
                <Icon name='home' />
                New Room
              </Button>
            </Menu.Item>
            <Menu.Item name='side layout' active >
              <Dropdown placeholder='Room' name="roomlist"  value={this.state.currentRoom} selection options={this.state.roomOptions} onChange={(event,{key,value})=>this.setState({currentRoom:value,currentRoomId:key})}/>
            </Menu.Item>
            <Menu.Item name='login' position="right" active >
              <Button icon labelPosition='left' color="teal" onClick={this.loginFormOpen}>
                <Icon name='user circle' />
                Login
              </Button>
            </Menu.Item>
          </Menu>
      </div>
      <Grid  divided className="bottom attached segment contentHeight">
        <Grid.Column width={3} stretched>
          <UserList userList={this.state.userList}/>
        </Grid.Column>
        <Grid.Column width={13} stretched className="contentHeight">
          <ChatMessageBox onClick={this.sendMessage} messageList={this.state.messageList}/>
        </Grid.Column>
      </Grid>
      <Segment basic></Segment>

      <Modal size="mini" open={this.state.modalOpen} onClose={this.close}>
          <Modal.Header>
            New User
          </Modal.Header>
          <Modal.Content>
          <Input fluid name="username" placeholder='username'  value={this.state.username}  onChange={this.setUsername}/>
          </Modal.Content>
          <Modal.Actions>
            <Button negative onClick={this.close}>
              No
            </Button>
            <Button positive icon='checkmark' labelPosition='right' content='Yes' onClick={this.logIn}/>
          </Modal.Actions>
        </Modal>
        <Modal size="mini" open={this.state.roomModal} onClose={this.close}>
          <Modal.Header>
            Add new room
          </Modal.Header>
          <Modal.Content>
          <Input fluid name="roomname" placeholder='room name'  value={this.state.newRoomName}  onChange={this.setNewRoomname}/>
          </Modal.Content>
          <Modal.Actions>
            <Button negative onClick={this.close}>
              No
            </Button>
            <Button positive icon='checkmark' labelPosition='right' content='Yes' onClick={this.createNewRoom}/>
          </Modal.Actions>
        </Modal>
    </Container>
    );
  }
}
export default Chatroom;
