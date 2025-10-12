import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Paper,
  Typography,
  Box,
  Button,
  AppBar,
  Toolbar,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Alert,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Tooltip,
  Collapse,
  LinearProgress,
  Grid,
} from '@mui/material';
import {
  Logout as LogoutIcon,
  PlayArrow as PlayArrowIcon,
  Stop as StopIcon,
  Refresh as RefreshIcon,
  Delete as DeleteIcon,
  Add as AddIcon,
  AccountCircle as AccountCircleIcon,
  Close as CloseIcon,
  AdminPanelSettings as AdminIcon,
  Info as InfoIcon,
  Article as ArticleIcon,
  BarChart as BarChartIcon,
  ExpandMore as ExpandMoreIcon,
} from '@mui/icons-material';
import { useAuth } from '../context/AuthContext';
import { containerApi } from '../services/api';

const DashboardPage = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [containers, setContainers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [newContainer, setNewContainer] = useState({
    name: '',
    image: '',
    envVars: []
  });
  const [envVarInput, setEnvVarInput] = useState('');
  const [actionLoading, setActionLoading] = useState({});
  const [logsDialogOpen, setLogsDialogOpen] = useState(false);
  const [selectedContainerId, setSelectedContainerId] = useState(null);
  const [containerLogs, setContainerLogs] = useState('');
  const [logsLoading, setLogsLoading] = useState(false);
  const [expandedRows, setExpandedRows] = useState({});
  const [containerStats, setContainerStats] = useState({});
  const [statsWebSockets, setStatsWebSockets] = useState({});
  const [logsWebSocket, setLogsWebSocket] = useState(null);
  const logsEndRef = React.useRef(null);

  useEffect(() => {
    if (user) {
      loadContainers();
    }
  }, [user]);

  useEffect(() => {
    return () => {
      Object.values(statsWebSockets).forEach(ws => {
        if (ws && ws.readyState === WebSocket.OPEN) {
          ws.close();
        }
      });
      if (logsWebSocket && logsWebSocket.readyState === WebSocket.OPEN) {
        logsWebSocket.close();
      }
    };
  }, [statsWebSockets, logsWebSocket]);

  const loadContainers = async () => {
    if (!user?.id) return;

    try {
      setLoading(true);
      setError('');
      const data = await containerApi.getAllContainers(user.id);
      setContainers(data || []);
    } catch (err) {
      setError('Failed to load containers: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const formatErrorMessage = (error) => {
    const message = error.response?.data?.message || error.message || 'Unknown error';

    if (message.includes('already in use') || message.includes('Conflict')) {
      const nameMatch = message.match(/name "\/([^"]+)"/);
      const containerName = nameMatch ? nameMatch[1] : 'this name';
      return `A container named "${containerName}" already exists. Please choose a different name or delete the existing container first.`;
    }

    if (message.includes('No such image') || message.includes('image not found')) {
      const imageMatch = message.match(/image[:\s]+([^\s:,]+)/i);
      const imageName = imageMatch ? imageMatch[1] : 'the image';
      return `Docker image "${imageName}" not found. Please check the image name and tag (e.g., nginx:latest).`;
    }

    if (message.includes('pull access denied') || message.includes('not found')) {
      return `Cannot pull the Docker image. Please verify the image name is correct and publicly accessible.`;
    }

    if (message.includes('Network') || message.includes('timeout') || message.includes('ECONNREFUSED')) {
      return `Cannot connect to Docker. Please ensure Docker is running.`;
    }

    if (message.includes('permission denied') || message.includes('access denied')) {
      return `Permission denied. Please check your Docker permissions.`;
    }

    const cleanMessage = message
      .replace(/Error response from daemon:\s*/i, '')
      .replace(/cannot create container:\s*/gi, '')
      .replace(/failed to create container:\s*/gi, '')
      .replace(/\s+/g, ' ')
      .trim();

    return cleanMessage.charAt(0).toUpperCase() + cleanMessage.slice(1);
  };

  const handleCreateContainer = async () => {
    try {
      setActionLoading({ create: true });
      setError('');

      await containerApi.createContainer({
        name: newContainer.name,
        image: newContainer.image,
        envVars: newContainer.envVars
      });

      setCreateDialogOpen(false);
      setNewContainer({ name: '', image: '', envVars: [] });
      setEnvVarInput('');
      await loadContainers();
    } catch (err) {
      setError(formatErrorMessage(err));
    } finally {
      setActionLoading({});
    }
  };

  const handleAddEnvVar = () => {
    if (envVarInput.trim()) {
      setNewContainer({
        ...newContainer,
        envVars: [...newContainer.envVars, envVarInput.trim()]
      });
      setEnvVarInput('');
    }
  };

  const handleRemoveEnvVar = (index) => {
    setNewContainer({
      ...newContainer,
      envVars: newContainer.envVars.filter((_, i) => i !== index)
    });
  };

  const handleStartContainer = async (containerId) => {
    try {
      setActionLoading({ [containerId]: 'start' });
      setError('');
      await containerApi.startContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to start container: ' + formatErrorMessage(err));
    } finally {
      setActionLoading({});
    }
  };

  const handleStopContainer = async (containerId) => {
    try {
      setActionLoading({ [containerId]: 'stop' });
      setError('');
      await containerApi.stopContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to stop container: ' + formatErrorMessage(err));
    } finally {
      setActionLoading({});
    }
  };

  const handleRestartContainer = async (containerId) => {
    try {
      setActionLoading({ [containerId]: 'restart' });
      setError('');
      await containerApi.restartContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to restart container: ' + formatErrorMessage(err));
    } finally {
      setActionLoading({});
    }
  };

  const handleDeleteContainer = async (containerId) => {
    if (!window.confirm('Are you sure you want to delete this container? This action cannot be undone.')) {
      return;
    }

    try {
      setActionLoading({ [containerId]: 'delete' });
      setError('');
      await containerApi.deleteContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to delete container: ' + formatErrorMessage(err));
    } finally {
      setActionLoading({});
    }
  };

  React.useEffect(() => {
    if (logsEndRef.current && containerLogs) {
      logsEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [containerLogs]);

  const handleViewLogs = (containerId) => {
    setSelectedContainerId(containerId);
    setLogsDialogOpen(true);
    setContainerLogs('Connecting to log stream...');
    setLogsLoading(true);

    if (logsWebSocket && logsWebSocket.readyState === WebSocket.OPEN) {
      logsWebSocket.close();
    }

    const wsUrl = containerApi.getContainerLogsWsUrl(containerId);


    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      setLogsLoading(false);
      setContainerLogs('');
    };

    ws.onmessage = (event) => {
      setContainerLogs(prev => prev + event.data);
    };

    ws.onerror = (error) => {
      setError(`Failed to connect to logs stream for container ${containerId.substring(0, 12)}`);
      setContainerLogs('Failed to connect to logs stream');
      setLogsLoading(false);
    };

    ws.onclose = (event) => {
      setLogsLoading(false);
      if (event.code !== 1000) {
        setError(`Logs WebSocket closed unexpectedly (code: ${event.code}${event.reason ? ', reason: ' + event.reason : ''})`);
      }
    };

    setLogsWebSocket(ws);
  };

  const handleToggleStats = (containerId) => {
    const isExpanded = expandedRows[containerId];

    if (isExpanded) {
      setExpandedRows(prev => ({ ...prev, [containerId]: false }));
      if (statsWebSockets[containerId]) {
        statsWebSockets[containerId].close();
        setStatsWebSockets(prev => {
          const newWs = { ...prev };
          delete newWs[containerId];
          return newWs;
        });
      }
      setContainerStats(prev => {
        const newStats = { ...prev };
        delete newStats[containerId];
        return newStats;
      });
    } else {
      setExpandedRows(prev => ({ ...prev, [containerId]: true }));
      connectStatsWebSocket(containerId);
    }
  };

  const connectStatsWebSocket = (containerId) => {
    const token = localStorage.getItem('token');
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/yadoma/ws/containers/${containerId}/stats?token=${encodeURIComponent(token)}`;
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {};

    ws.onmessage = (event) => {
      try {
        const stats = JSON.parse(event.data);

        if (stats.error) {
          setError('Stats error: ' + stats.error);
          setExpandedRows(prev => ({ ...prev, [containerId]: false }));
          ws.close();
          return;
        }

        const formattedStats = {
          ...stats,
          cpu: (stats.cpu / 1000000000).toFixed(2)
        };

        setContainerStats(prev => ({ ...prev, [containerId]: formattedStats }));
      } catch (err) {
        setError('Failed to parse stats data');
      }
    };

    ws.onerror = (error) => {
      setError('Failed to connect to stats stream. Make sure the container is running.');
      setExpandedRows(prev => ({ ...prev, [containerId]: false }));
    };

    ws.onclose = (event) => {
      if (event.code !== 1000) {
        setError(`Stats WebSocket closed unexpectedly (code: ${event.code}${event.reason ? ', reason: ' + event.reason : ''})`);
      }
    };

    setStatsWebSockets(prev => ({ ...prev, [containerId]: ws }));
  };

  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getStatusColor = (status) => {
    const statusLower = status?.toLowerCase();
    if (statusLower?.includes('up')) {
      return 'success';
    }
    switch (statusLower) {
      case 'running':
        return 'success';
      case 'stopped':
      case 'exited':
        return 'error';
      case 'paused':
        return 'warning';
      default:
        return 'default';
    }
  };

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Yadoma - Container Management
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Button
              color="inherit"
              startIcon={<InfoIcon />}
              onClick={() => navigate('/system')}
            >
              System Info
            </Button>
            {user?.role === 'ADMIN' && (
              <Button
                color="inherit"
                startIcon={<AdminIcon />}
                onClick={() => navigate('/admin/users')}
              >
                Admin Panel
              </Button>
            )}
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <AccountCircleIcon />
              <Typography variant="body1">{user?.email}</Typography>
              {user?.role === 'ADMIN' && (
                <Chip label="ADMIN" color="error" size="small" />
              )}
            </Box>
            <IconButton color="inherit" onClick={handleLogout}>
              <LogoutIcon />
            </IconButton>
          </Box>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Paper sx={{ p: 3 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h5" component="h2">
              My Containers
            </Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => setCreateDialogOpen(true)}
            >
              Create Container
            </Button>
          </Box>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
              {error}
            </Alert>
          )}

          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : containers.length === 0 ? (
            <Box sx={{ textAlign: 'center', py: 4 }}>
              <Typography variant="body1" color="text.secondary">
                No containers found. Create your first container to get started.
              </Typography>
            </Box>
          ) : (
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>State</TableCell>
                    <TableCell align="right">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {containers.map((container) => {
                    const isRunning = container.state?.toLowerCase() === 'running' ||
                                     container.status?.toLowerCase()?.includes('up');
                    const isStopped = !isRunning;

                    return (
                      <React.Fragment key={container.id}>
                        <TableRow>
                          <TableCell>{container.name}</TableCell>
                          <TableCell>
                            <Chip
                              label={container.status}
                              color={getStatusColor(container.status)}
                              size="small"
                            />
                          </TableCell>
                          <TableCell>{container.state}</TableCell>
                          <TableCell align="right">
                            <Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
                              <Tooltip title={isRunning ? "View Stats" : "Container must be running"}>
                                <span>
                                  <IconButton
                                    color="primary"
                                    size="small"
                                    onClick={() => handleToggleStats(container.id)}
                                    disabled={!isRunning}
                                    sx={{ opacity: isRunning ? 1 : 0.3 }}
                                  >
                                    {expandedRows[container.id] ? <ExpandMoreIcon sx={{ transform: 'rotate(180deg)' }} /> : <BarChartIcon />}
                                  </IconButton>
                                </span>
                              </Tooltip>
                              <Tooltip title={isRunning ? "View Logs" : "Container must be running"}>
                                <span>
                                  <IconButton
                                    color="info"
                                    size="small"
                                    onClick={() => handleViewLogs(container.id)}
                                    disabled={!isRunning}
                                    sx={{ opacity: isRunning ? 1 : 0.3 }}
                                  >
                                    <ArticleIcon />
                                  </IconButton>
                                </span>
                              </Tooltip>
                          <Tooltip title="Start">
                            <span>
                              <IconButton
                                color="success"
                                size="small"
                                onClick={() => handleStartContainer(container.id)}
                                disabled={isRunning || !!actionLoading[container.id]}
                                sx={{ opacity: isRunning ? 0.3 : 1 }}
                              >
                                {actionLoading[container.id] === 'start' ? (
                                  <CircularProgress size={20} />
                                ) : (
                                  <PlayArrowIcon />
                                )}
                              </IconButton>
                            </span>
                          </Tooltip>
                          <Tooltip title="Stop">
                            <span>
                              <IconButton
                                color="error"
                                size="small"
                                onClick={() => handleStopContainer(container.id)}
                                disabled={isStopped || !!actionLoading[container.id]}
                                sx={{ opacity: isStopped ? 0.3 : 1 }}
                              >
                                {actionLoading[container.id] === 'stop' ? (
                                  <CircularProgress size={20} />
                                ) : (
                                  <StopIcon />
                                )}
                              </IconButton>
                            </span>
                          </Tooltip>
                          <Tooltip title="Restart">
                            <span>
                              <IconButton
                                color="primary"
                                size="small"
                                onClick={() => handleRestartContainer(container.id)}
                                disabled={isStopped || !!actionLoading[container.id]}
                                sx={{ opacity: isStopped ? 0.3 : 1 }}
                              >
                                {actionLoading[container.id] === 'restart' ? (
                                  <CircularProgress size={20} />
                                ) : (
                                  <RefreshIcon />
                                )}
                              </IconButton>
                            </span>
                          </Tooltip>
                          <Tooltip title="Delete">
                            <span>
                              <IconButton
                                color="error"
                                size="small"
                                onClick={() => handleDeleteContainer(container.id)}
                                disabled={isRunning || !!actionLoading[container.id]}
                                sx={{ opacity: isRunning ? 0.3 : 1 }}
                              >
                                {actionLoading[container.id] === 'delete' ? (
                                  <CircularProgress size={20} />
                                ) : (
                                  <DeleteIcon />
                                )}
                              </IconButton>
                            </span>
                          </Tooltip>
                        </Box>
                      </TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={4}>
                        <Collapse in={expandedRows[container.id]} timeout="auto" unmountOnExit>
                          <Box sx={{ margin: 2 }}>
                            <Typography variant="h6" gutterBottom component="div">
                              Real-time Statistics
                            </Typography>
                            {containerStats[container.id] ? (
                              <Grid container spacing={2}>
                                <Grid item xs={12} md={6}>
                                  <Paper sx={{ p: 2 }}>
                                    <Typography variant="subtitle2" color="text.secondary">
                                      CPU Usage
                                    </Typography>
                                    <Typography variant="h4">
                                      {containerStats[container.id].cpu}%
                                    </Typography>
                                    <LinearProgress
                                      variant="determinate"
                                      value={containerStats[container.id].cpu}
                                      sx={{ mt: 1 }}
                                    />
                                  </Paper>
                                </Grid>
                                <Grid item xs={12} md={6}>
                                  <Paper sx={{ p: 2 }}>
                                    <Typography variant="subtitle2" color="text.secondary">
                                      Memory Usage
                                    </Typography>
                                    <Typography variant="h6">
                                      {formatBytes(containerStats[container.id].memUsage)} / {formatBytes(containerStats[container.id].memLimit)}
                                    </Typography>
                                    <LinearProgress
                                      variant="determinate"
                                      value={(containerStats[container.id].memUsage / containerStats[container.id].memLimit) * 100}
                                      sx={{ mt: 1 }}
                                      color="secondary"
                                    />
                                  </Paper>
                                </Grid>
                                <Grid item xs={12} md={6}>
                                  <Paper sx={{ p: 2 }}>
                                    <Typography variant="subtitle2" color="text.secondary">
                                      Network Input
                                    </Typography>
                                    <Typography variant="h6">
                                      {formatBytes(containerStats[container.id].netInput)}
                                    </Typography>
                                  </Paper>
                                </Grid>
                                <Grid item xs={12} md={6}>
                                  <Paper sx={{ p: 2 }}>
                                    <Typography variant="subtitle2" color="text.secondary">
                                      Network Output
                                    </Typography>
                                    <Typography variant="h6">
                                      {formatBytes(containerStats[container.id].netOutput)}
                                    </Typography>
                                  </Paper>
                                </Grid>
                              </Grid>
                            ) : (
                              <Box sx={{ display: 'flex', justifyContent: 'center', py: 2 }}>
                                <CircularProgress size={24} />
                                <Typography sx={{ ml: 2 }}>Loading stats...</Typography>
                              </Box>
                            )}
                          </Box>
                        </Collapse>
                      </TableCell>
                    </TableRow>
                  </React.Fragment>
                    );
                  })}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Paper>
      </Container>

      <Dialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Create New Container</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            <TextField
              autoFocus
              label="Container Name"
              fullWidth
              variant="outlined"
              required
              value={newContainer.name}
              onChange={(e) => setNewContainer({ ...newContainer, name: e.target.value.trim() })}
              helperText="A unique name for your container (e.g., my-nginx, redis-cache)"
              error={newContainer.name && !/^[a-zA-Z0-9][a-zA-Z0-9_.-]*$/.test(newContainer.name)}
            />
            <TextField
              label="Docker Image"
              fullWidth
              variant="outlined"
              required
              value={newContainer.image}
              onChange={(e) => setNewContainer({ ...newContainer, image: e.target.value.trim() })}
              helperText="e.g., nginx:latest, ubuntu:22.04, postgres:15, redis:alpine"
              placeholder="nginx:latest"
            />

            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Environment Variables (optional)
              </Typography>
              <Box sx={{ display: 'flex', gap: 1, mb: 1 }}>
                <TextField
                  size="small"
                  fullWidth
                  variant="outlined"
                  placeholder="KEY=value"
                  value={envVarInput}
                  onChange={(e) => setEnvVarInput(e.target.value)}
                  onKeyPress={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      handleAddEnvVar();
                    }
                  }}
                />
                <Button
                  variant="outlined"
                  onClick={handleAddEnvVar}
                  disabled={!envVarInput.trim()}
                >
                  Add
                </Button>
              </Box>

              {newContainer.envVars.length > 0 && (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5 }}>
                  {newContainer.envVars.map((envVar, index) => (
                    <Chip
                      key={index}
                      label={envVar}
                      onDelete={() => handleRemoveEnvVar(index)}
                      deleteIcon={<CloseIcon />}
                      variant="outlined"
                      size="small"
                    />
                  ))}
                </Box>
              )}
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => {
            setCreateDialogOpen(false);
            setNewContainer({ name: '', image: '', envVars: [] });
            setEnvVarInput('');
          }}>
            Cancel
          </Button>
          <Button
            onClick={handleCreateContainer}
            variant="contained"
            disabled={
              !newContainer.name ||
              !newContainer.image ||
              actionLoading.create ||
              (newContainer.name && !/^[a-zA-Z0-9][a-zA-Z0-9_.-]*$/.test(newContainer.name))
            }
          >
            {actionLoading.create ? <CircularProgress size={24} /> : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog
        open={logsDialogOpen}
        onClose={() => {
          if (logsWebSocket && logsWebSocket.readyState === WebSocket.OPEN) {
            logsWebSocket.close();
          }
          setLogsWebSocket(null);
          setLogsDialogOpen(false);
          setContainerLogs('');
        }}
        maxWidth="lg"
        fullWidth
      >
        <DialogTitle>
          Container Logs
          {logsLoading && <CircularProgress size={20} sx={{ ml: 2 }} />}
        </DialogTitle>
        <DialogContent>
          <Box
            sx={{
              mt: 1,
              p: 2,
              backgroundColor: '#1e1e1e',
              color: '#d4d4d4',
              fontFamily: 'Consolas, Monaco, "Courier New", monospace',
              fontSize: '13px',
              borderRadius: '4px',
              border: '1px solid #3e3e3e',
              maxHeight: '500px',
              overflow: 'auto',
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-word',
            }}
          >
            {containerLogs}
            <div ref={logsEndRef} />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => {
            if (logsWebSocket && logsWebSocket.readyState === WebSocket.OPEN) {
              logsWebSocket.close();
            }
            setLogsWebSocket(null);
            setLogsDialogOpen(false);
            setContainerLogs('');
          }}>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default DashboardPage;
