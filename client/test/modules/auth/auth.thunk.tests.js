import { spy } from 'sinon';
import { expect } from 'chai';
import auth from '../../../src/modules/auth/auth.thunk';

describe('I wish to verify the authentification, register, resetPassword and changePasswordAfterReset are using the method POST', () => {

  it('The method should be equal to POST', () => {

    // Given
    let obj = {};

    // When
    let config = auth.encodeInfos(obj);

    // Then
    expect(config.method).to.be.equal('POST');
  });
});

describe('I wish to verify that credentials are passed correctly when I log', () => {

  // Given
  let creds = {
    username: 'username',
    password: 'password'
  };

  it('I verify that infos used for authentification are used by encodeURIComponent', () => {

    // Spy
    let encodeURIComponentSpy = spy(window, 'encodeURIComponent');

    // When
    auth.encodeInfos(creds);

    // Then
    expect(encodeURIComponentSpy.callCount).to.be.equal(Object.keys(creds).length * 2);

    // Clean
    encodeURIComponentSpy.restore();
  });

  it('Hashs of username and password should be equal to those generated within encodeInfos (nominal case)', () => {

    // Given
    let encodedUsername = window.encodeURIComponent(creds.username);
    let encodedPassword = window.encodeURIComponent(creds.password);

    // When
    let config = auth.encodeInfos(creds);

    // Then
    expect(config.body).to.include('username=' + encodedUsername);
    expect(config.body).to.include('password=' + encodedPassword);
  });

  it('Hashs of username and password should be equal to those generated within encodeInfos (edge case)', () => {

    // Given
    let edgeCaseCreds = {
      username: '!:;@*&,',
      password: '!:;@*&!'
    };

    let encodedUsername = window.encodeURIComponent(edgeCaseCreds.username);
    let encodedPassword = window.encodeURIComponent(edgeCaseCreds.password);

    // When
    let config = auth.encodeInfos(edgeCaseCreds);

    // Then
    expect(config.body).to.include('username=' + encodedUsername);
    expect(config.body).to.include('password=' + encodedPassword);
  });
});

describe('I wish to verify that account\'s infos are passed correctly when I register', () => {

  // Given
  let account = {
    username:   'username',
    password:   'password',
    mail:       'mail@example.com',
    firstname:  'firsname',
    lastname:   'lastname'
  };

  it('I verify that account\'s infos used for register are used by encodeURIComponent', () => {

    // Spy
    let encodeURIComponentSpy = spy(window, 'encodeURIComponent');

    // When
    auth.encodeInfos(account);

    // Then
    expect(encodeURIComponentSpy.callCount).to.be.equal(Object.keys(account).length * 2);

    // Clean
    encodeURIComponentSpy.restore();
  });

  it('Hashs of account\'infos should be equal to those generated within encodeInfos', () => {

    // Given
    let encodedUsername = window.encodeURIComponent(account.username);
    let encodedPassword = window.encodeURIComponent(account.password);
    let encodedMail = window.encodeURIComponent(account.mail);
    let encodedFirstName = window.encodeURIComponent(account.firstname);
    let encodedLastName = window.encodeURIComponent(account.lastname);

    // When
    let config = auth.encodeInfos(account);

    // Then
    expect(config.body).to.include('username=' + encodedUsername);
    expect(config.body).to.include('password=' + encodedPassword);
    expect(config.body).to.include('mail=' + encodedMail);
    expect(config.body).to.include('firstname=' + encodedFirstName);
    expect(config.body).to.include('lastname=' + encodedLastName);
  });
});

describe('I wish to verify that username is passed correctly when I reset the password', () => {

  // Given
  let userNameObject = {
    username: 'username'
  };

  it('I verify the userNameObject used for the reset is used by encodeURIComponent', () => {

    // Spy
    let encodeURIComponentSpy = spy(window, 'encodeURIComponent');

    // When
    auth.encodeInfos(userNameObject);

    // Then
    expect(encodeURIComponentSpy.callCount).to.be.equal(Object.keys(userNameObject).length * 2);

    // Clean
    encodeURIComponentSpy.restore();
  });

  it('Hash of username should be equal to the one generated within encodeInfos', () => {

    // Given
    let encodedUsername = window.encodeURIComponent(userNameObject.username);

    // When
    let config = auth.encodeInfos(userNameObject);

    // Then
    expect(config.body).to.include(encodedUsername);
  });
});

describe('I wish to verify that username is passed correctly when I change the password after a reset', () => {

  // Given
  let newPwdTokenObject = {
    newPassword:  'newPassword',
    token:        'token'
  };

  it('I verify the newPwdTokenObject used for changing the password is used by encodeURIComponent', () => {

    // Spy
    let encodeURIComponentSpy = spy(window, 'encodeURIComponent');

    // When
    auth.encodeInfos(newPwdTokenObject);

    // Then
    expect(encodeURIComponentSpy.callCount).to.be.equal(Object.keys(newPwdTokenObject).length * 2);

    // Clean
    encodeURIComponentSpy.restore();
  });

  it('Hash of username should be equal to the one generated within encodeInfos', () => {

    // Given
    let encodedPassword = window.encodeURIComponent(newPwdTokenObject.newPassword);
    let encodedToken = window.encodeURIComponent(newPwdTokenObject.token);

    // When
    let config = auth.encodeInfos(newPwdTokenObject);

    // Then
    expect(config.body).to.include(encodedPassword);
    expect(config.body).to.include(encodedToken);
  });
});
