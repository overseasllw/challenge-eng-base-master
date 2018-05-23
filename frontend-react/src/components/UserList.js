import React, { Component } from 'react';
import { List } from 'semantic-ui-react'
class UserList extends Component{
    render(){
        return (
            <div>
              <List as="menu" divided relaxed>
                <List.Header>Users</List.Header>
                <List.Item>
                  <List.Icon name='user' size='large' verticalAlign='middle' />
                  <List.Content>
                    <List.Header as='a'>liwei</List.Header>
                    <List.Description as='a'>Updated 10 mins ago</List.Description>
                  </List.Content>
                </List.Item>
                <List.Item>
                  <List.Icon name='user' size='large' verticalAlign='middle' />
                  <List.Content>
                    <List.Header as='a'>john</List.Header>
                    <List.Description as='a'>Updated 22 mins ago</List.Description>
                  </List.Content>
                </List.Item>
                <List.Item>
                  <List.Icon name='user' size='large' verticalAlign='middle' />
                  <List.Content>
                    <List.Header as='a'>mike</List.Header>
                    <List.Description as='a'>Updated 34 mins ago</List.Description>
                  </List.Content>
                </List.Item>
              </List>
            </div>
          );
    }
}

export default UserList;