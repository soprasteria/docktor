import React from 'react'
import { IndexLink, Link } from 'react-router';

import 'semantic-ui-menu/menu.min.css'
import 'semantic-ui-icon/icon.min.css'
import './navBar.scss'



const NavBar = () => (
    <div className="ui inverted fluid menu navbar">
      <IndexLink to="/" className="item">
        <i className="big fitted doctor icon"></i>{' '}Docktor
      </IndexLink>
      <Link to="/sites" activeClassName="active" className="item">Sites</Link>
      <Link to="/daemons" activeClassName="active" className="item">Daemons</Link>
      <Link to="/groups" activeClassName="active" className="item">Groups</Link>
      <Link to="/users" activeClassName="active" className="item">Users</Link>
      <div className="right menu">
        <a href="#" className="item"><i className="inverted large user icon"></i>{' '} Admin</a>
      </div>
    </div>
  )

export default NavBar;