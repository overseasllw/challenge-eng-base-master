import React, { Component } from 'react';
import Chatroom from './components/Chatroom'

class App extends Component {
  render() {
    return (
      <Chatroom/>
    );
  }
  /*componentDidMount() {
    fetch('/api/v1/messages/').then((res) => {
      return res.json();
    }).then((res) => {
      this.setState({res});
    }).catch((err) => {
      this.setState({err});
    });
  }*/
}

export default App;
