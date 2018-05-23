import React, { Component } from 'react';
import { Grid, Container,Menu,Button, Icon, Segment, Modal,Input} from 'semantic-ui-react'
import UserList from './UserList';
import ChatMessageBox from './ChatMessageBox'

class Chatroom extends Component{
  messages =[]
  constructor(props){
    super(props)
    this.state = {
      message:"",
      //:[],
      username:"guest_"+this.makeid(),
      modalOpen:false,
    }
   // this.ws //= this.ws.bind(this);
    this.initSocket = this.initSocket.bind(this);
    this.sendMessage = this.sendMessage.bind(this);
    this.loginFormOpen =this.loginFormOpen.bind(this)
    this.setUsername = this.setUsername.bind(this)
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
    this.ws = new WebSocket("ws://localhost:18000/ws");
    this.ws.onmessage = (msg) => {
      //console.log(msg.data)
      this.setState({ message: msg.data });
   //   this.state.messages.concat(this.state.messages);
      console.log(this.state.message)
    }
    console.log("init")
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
     // this.setState({messages:res})
     this.messages = res
      console.log(this.messages)
    }).catch((err) => {
      this.setState({err});
    });
  }

  generateTimestamp () {
    var iso = new Date().toISOString();
    return iso.split("T")[1].split(".")[0];
  }

  sendMessage (message) {
      this.ws.send(
        JSON.stringify({
          username: this.state.username,
          message: (this.generateTimestamp() + " <"+ this.state.username +"> " + message)
        })
      );
  }

  loginFormOpen(event){
    event.preventDefault();
    this.setState(prevState => ({...prevState, modalOpen: true}) )
  }
  close = () => this.setState(prevState =>({...prevState, modalOpen: false }))
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
          <UserList/>
        </Grid.Column>
        <Grid.Column width={13} stretched className="contentHeight">
          <ChatMessageBox onClick={this.sendMessage} messages={this.massages}/>
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
            <Button positive icon='checkmark' labelPosition='right' content='Yes' onClick={this.sendMessage}/>
          </Modal.Actions>
        </Modal>
    </Container>
    );
  }
}
export default Chatroom;