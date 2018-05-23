import React, { Component } from 'react';
import { Button,Icon, Grid, Input } from 'semantic-ui-react'
class ChatInputBox extends Component{
    constructor(props){
        super(props)
        this.state = {message:''}
        this.messageOnKeyDown = this.messageOnKeyDown.bind(this);
        this.sendButtononClick = this.sendButtononClick.bind(this)
        this.messageOnChange = this.messageOnChange.bind(this)
    }
    render(){
        return (
            <div className="inputBox">
                <Grid stackable>
                    <Grid.Column width={13}>
                        <Input fluid name="mes" placeholder='message ...'  value={this.state.message} onKeyDown={this.messageOnKeyDown} onChange={this.messageOnChange}/>
                    </Grid.Column>
                    <Grid.Column width={3}>
                        <Button icon color="green"  fluid onClick={this.sendButtononClick}>
                        <Icon name="send"/>
                            Send
                        </Button>
                    </Grid.Column>
                </Grid>
            </div>
        )
    }

    sendButtononClick(event) {
        this.props.onClick(this.state.message)
        this.setState({
            message: ''
        });
    }

    messageOnChange(event){
        this.setState({message:event.target.value})
    }

    messageOnKeyDown(event) {
        if (event.keyCode === 13 && this.state.message !== "") {
         // this.props.onSubmit(this.state.message);
          //this.state.messages.push(event.target.value)
       //   console.log(this.state.messages)
       //   console.log(this.state.message)
          this.props.onClick(this.state.message)
          this.setState({ message: '' });
        }
      }
    
}

export default ChatInputBox;