import React, { Component } from 'react';
import {  Container,Divider,Label,Segment } from 'semantic-ui-react'
import ChatInputBox from './ChatInputBox'
class ChatMessageBox extends Component{
   /* constructor(props){
        super(props)
    }*/
    render(){
        return (
            <Segment className="contentHeight">
                <Container className="messageBox transition visible"  >
                {
                    this.props.messageList.map((m)=> 
                    {
                        if(m.message_type!=="system-message"){
                            return  <div className="message" key={m.uuid}>
                                    <Label>{m.username}</Label>
                                    <Label color='blue' pointing='left'>{m.message}</Label>
                                </div>
                        }
                            return <div className="message" key={m.uuid}>
                                    <Label color='orange'>{m.message}</Label>
                                </div>
                        
                    })    
                }
                </Container>
                <Divider horizontal>Message</Divider>
                <ChatInputBox onClick={this.props.onClick} typingIndicator={this.props.typingIndicator} removeIndicator={this.props.removeIndicator}/>
            </Segment>
        );
    }
    /*shouldComponentUpdate(nextProps){
        return !(nextProps.messageList!== this.props.messageList)
    }*/
    
}

export default ChatMessageBox