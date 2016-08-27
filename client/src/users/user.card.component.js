// React
import React from 'react';
import { Link } from 'react-router';

// Style
import './user.card.component.scss';

// UserCard Component
class UserCard extends React.Component {
    render() {
      const user = this.props.user;
        return (
          <div className='ui card user'>
            <div className='content'>
              <div className='right floated meta'>{user.Role}</div>
              <img className='ui avatar image' src='./images/avatar.jpg'/>{user.DisplayName}
            </div>
            <div className='extra content'>
              <span className='right floated'>
                <i className='travel icon'></i>
                {user.Groups.length + ' group(s)'}
              </span>
              <div className='email' title={user.Email}>
              <i className='mail icon'></i>{user.Email}
              </div>
            </div>
          </div>
        );
    }
}
UserCard.propTypes = { user: React.PropTypes.object };

export default UserCard;