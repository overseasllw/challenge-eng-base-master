import React, { Component } from 'react';
import {  Container,Divider,Label,Segment } from 'semantic-ui-react'
import ChatInputBox from './ChatInputBox'
class ChatMessageBox extends Component{
   /* constructor(props){
        super(props)
     //   console.log(this.props.messages)
        this.messages = []//this.props.messages
    }*/
    render(){
        var ms = this.props.messages.map(function(m,index) {
           return( <div className="message" key={m.message_id}>
                <Label>{m.username}</Label>
                <Label color='blue' pointing='left'>{m.message}</Label>
            </div>)
            })
        return (
            <Segment className="contentHeight">
                <Container className="messageBox">
                {ms}
                </Container>
                <Divider horizontal>Message</Divider>
                <ChatInputBox onClick={this.props.onClick}/>
            </Segment>
        );
    }
    shouldComponentUpdate(nextProps){
        return !(nextProps.messages!== this.props.messages)
    }
    
}

export default ChatMessageBox