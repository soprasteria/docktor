
import groupBy from 'lodash.groupby';
import sortBy from 'lodash.sortby';

import { transformFilterToObject } from '../utils/search.js';
import { containsWithoutAccents } from '../utils/strings.js';

export const getFilteredGroups = (groups, filterValue) => {
  if (!filterValue || filterValue === '') {
    return Object.values(groups);
  } else {
    return Object.values(groups).filter(group => {
      let match = true;
      const query = transformFilterToObject(filterValue);
      Object.keys(query).forEach(key => {
        const value = query[key];
        switch(key) {
        case 'text':
          match &= containsWithoutAccents(JSON.stringify(Object.values(group)), value);
          return;
        case 'name':
        case 'title':
          match &= containsWithoutAccents(group.title, value);
          return;
        case 'tags':
          const tags = group.tags || [];
          match &= tags.filter(tag => containsWithoutAccents(tag, value)).length > 0;
          return;
        default:
          match = false;
          return;
        }
      });
      return match;
    });
  }
};

export const getContainersGroupByCategory = (containers, services, tags) => {

  if (tags.isFetching || services.isFetching || !tags.items || !services.items || Object.keys(tags.items) == 0 || Object.keys(services.items) == 0) {
    return [];
  }

  const enrichedContainers = containers.map(container => {
    const serviceTags = services.items[container.serviceId] ? (services.items[container.serviceId].tags || []) : [];
    const pack = serviceTags.map(tag => tags.items[tag]).filter(tag => tag.category.slug === 'package').map(tag => tag.name.raw);
    container.package = pack[0] || 'Others';
    return container;
  });

  const sortedContainers = sortBy(enrichedContainers, t => t.title);
  const groupByPackage = groupBy(sortedContainers, (c => c.package));
  return groupByPackage;
};
