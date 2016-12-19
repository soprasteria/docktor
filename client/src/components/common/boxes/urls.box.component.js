// React
import React from 'react';

import Box from './box/box.component.js';

// URLsBox is a list of docker urls
class URLsBox extends React.Component {

  constructor(props) {
    super(props);

    // Set state of component from the props.
    this.state = { urls: this.props.urls || [] };
  }

  isFormValid() {
    return this.refs.urlsBox.isFormValid();
  }

  onChangeURLs(urls) {
    this.state.urls = urls;
  }

  render() {
    const form = { fields:[] };

    form.getTitle = (url) => {
      return '';
    };

    form.fields.push({
      name: 'label',
      label: 'Label',
      placeholder: 'Label',
      sizeClass: 'seven wide',
      isRequired: true
    });

    form.fields.push({
      name: 'url',
      label: 'URL',
      placeholder: 'URL',
      sizeClass: 'eight wide',
      isRequired: true
    });

    return (
      <Box
        ref='urlsBox'
        boxId={this.props.boxId}
        icon='large linkify icon'
        title='URLs' form={form}
        lines={this.props.urls}
        stacked={this.props.stacked}
        onChange={urls => this.onChangeURLs(urls)}>
        {this.props.children || ''}
      </Box>
    );
  }
}

URLsBox.propTypes = {
  boxId: React.PropTypes.string,
  urls: React.PropTypes.array,
  stacked: React.PropTypes.bool,
  children: React.PropTypes.oneOfType([
    React.PropTypes.array,
    React.PropTypes.element
  ])
};

export default URLsBox;