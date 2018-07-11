import React, { Component } from 'react';
import { List } from 'semantic-ui-react'
class UserList extends Component{
    render(){
        return (
            <div>
              <List as="menu" divided relaxed>
                <List.Header>Users</List.Header>
                {
                  this.props.userList.map((user)=>
                  <List.Item key={user.username}>
                    <List.Icon name='user' size='large' verticalAlign='middle' />
                    <List.Content>
                      <List.Header as='a'>{user.username}</List.Header>
                      <List.Description as='a'>*</List.Description>
                    </List.Content>
                  </List.Item>)
                }
              </List>
            </div>
          );
    }
}

export default UserList;
