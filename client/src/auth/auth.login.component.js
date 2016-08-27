// React
import React from 'react';
import classNames from 'classnames';

// Style
import './auth.component.scss';

// Signin Pane containing fields to log in the application
class SigninPane extends React.Component {

  componentDidMount() {
    $('.ui.form')
      .form({
        fields: {
          username : 'empty',
          password : 'empty',
        },
        onSuccess: (event, fields) => {
          this.handleClick(event);
        },
        onFailure: (event, fields) => {
          return false;
        }
      })
    ;
  }

  render() {
    const { errorMessage } = this.props;
    return (
      <div id='login'>
        <h1>{this.props.title}</h1>
        <form className='ui form'>
          <div className='field'>
            <label>
              Username<span className='req'>*</span>
            </label>
            <input type='text' ref='username' name='username' autoComplete='off' placeholder='Registered/LDAP username' />
          </div>
          <div className='field'>
              <label>
              Password<span className='req'>*</span>
              </label>
              <input type='password' ref='password' name='password' autoComplete='off' placeholder='Password' />
          </div>
          {errorMessage &&
              <p className='error'>{errorMessage}</p>
          }
          <div className='ui error message'></div>
          <p className='forgot'><a href='#'>Forgot Password?</a></p>
          <button type='submit' className='button button-block'>{this.props.submit}</button>
        </form>
      </div>
    );
  }

  handleClick(event) {
      event.preventDefault();
      const username = this.refs.username;
      const password = this.refs.password;
      const creds = { username: username.value.trim(), password: password.value.trim() };
      this.props.onLoginClick(creds);
  }
};

SigninPane.propTypes = {
  onLoginClick: React.PropTypes.func.isRequired,
  errorMessage: React.PropTypes.string,
  label: React.PropTypes.string.isRequired,
  title: React.PropTypes.string.isRequired,
  submit: React.PropTypes.string.isRequired
};

export default SigninPane;