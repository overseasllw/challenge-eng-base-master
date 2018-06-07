import React, { Component } from 'react';
import { Grid, Container,Menu,Button, Icon, Segment, Modal,Input} from 'semantic-ui-react'
import UserList from './UserList';
import ChatMessageBox from './ChatMessageBox'

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
      modalOpen:false,
    }
   // this.setState({messageList:[]})
    this.initSocket = this.initSocket.bind(this);
    this.sendMessage = this.sendMessage.bind(this);
    this.loginFormOpen =this.loginFormOpen.bind(this)
    this.setUsername = this.setUsername.bind(this)
    this.close = this.close.bind(this)
    this.logIn = this.logIn.bind(this)
  }

   makeid() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  
    for (var i = 0; i < 6; i++)
      text += possible.charAt(Math.floor(Math.random() * possible.length));
  
    return text;
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
    fetch("/api/v1/messages/").then((res) => {
      return res.json();
    }).then((res) => {
       // console.log(this.state.messages)
      this.setState({messageList:res})
      this.messageList = res
    //  console.log(this.messageList)
    }).catch((err) => {
      this.setState({err});
    });
  }

  generateTimestamp () {
    var iso = new Date().toTimeString() //.toISOString();
    return iso;//iso.split("T")[1].split(".")[0];
  }

  sendMessage (message) {
      this.ws.send(
        JSON.stringify({
          username: this.state.username===""?this.state.guestname:this.state.username,
          register:false,
          guestname:this.state.guestname,
          message: (message+" ("+this.generateTimestamp()+")")
        })
      );
  }
  logIn(){
    this.ws.send(
      JSON.stringify({
        username:this.state.username,
        guestname:this.state.guestname,
        register:true,
        message:""
      })
    )
    this.setState({modalOpen:false})
  }

  loginFormOpen(event){
    event.preventDefault();
    this.setState({ modalOpen: true})
  }
  close(event){
    this.setState({ modalOpen: false })
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
            Delete Your Account
          </Modal.Header>
          <Modal.Content>
          <Input fluid name="username" placeholder='username'  value={this.state.username}  onChange={this.setUsername}/>
          </Modal.Content>
          <Modal.Actions>
            <Button negative>
              No
            </Button>
            <Button positive icon='checkmark' labelPosition='right' content='Yes' onClick={this.logIn}/>
          </Modal.Actions>
        </Modal>
    </Container>
    );
  }
}
export default Chatroom;