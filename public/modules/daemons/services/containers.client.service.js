'use strict';

angular.module('daemons').factory('Containers', ['$http', 'DaemonsDocker',
    function ($http, DaemonsDocker) {
        return {
            inspectContainer: function (daemonId, containerId) {
                return $http.get('/daemons/docker/container/inspect/' + daemonId + '/' + containerId);
            },
            statsContainer: function (container, machineInfo, daemonId, containerId) {
                return $http.get('/daemons/docker/container/stats/' + daemonId + '/' + containerId).
                    success(function (containerInfo, status, headers, config) {
                        container.statsCompute = containerInfo.stats;

                        var cur = containerInfo.stats[containerInfo.stats.length - 1];
                        var cpuUsage = 0;
                        if (containerInfo.spec.has_cpu && containerInfo.stats.length >= 2) {
                            var prev = containerInfo.stats[containerInfo.stats.length - 2];
                            var rawUsage = cur.cpu.usage.total - prev.cpu.usage.total;
                            var intervalInNs = DaemonsDocker.getInterval(cur.timestamp, prev.timestamp);
                            // Convert to millicores and take the percentage
                            cpuUsage = Math.round(((rawUsage / intervalInNs) / machineInfo.num_cores) * 100);
                            if (cpuUsage > 100) {
                                cpuUsage = 100;
                            }
                        }
                        container.statsCompute.cpuUsagePercent = cpuUsage;

                        if (containerInfo.spec.has_memory) {
                            // Saturate to the machine size.
                            var limit = containerInfo.spec.memory.limit;
                            if (limit > machineInfo.memory_capacity) {
                                limit = machineInfo.memory_capacity;
                            }
                            container.statsCompute.memoryLimit = limit;
                            container.statsCompute.memoryUsage = Math.round(cur.memory.usage / 1000000);
                            container.statsCompute.memoryUsagePercent = Math.round((cur.memory.usage / limit) * 100);
                        }
                    }).
                    error(function (data, status, headers, config) {
                        console.log('Error:');
                        console.log(data);
                    });
            },
            actionContainer: function (action, daemonId, container, cbSuccess, index, cbSuccessEnd, cbError) {
                return $http.get('/daemons/docker/container/' + action + '/' + daemonId + '/' + container.Id).
                    success(function (data, status, headers, config) {
                        cbSuccess(container, data, index, cbSuccessEnd);
                    }).
                    error(function (data, status, headers, config) {
                        cbError(container, data, index);
                    });
            }
        };
    }
]);
