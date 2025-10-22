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

  useEffect(() => {
    if (user) {
      loadContainers();
    }
  }, [user]);

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

  const handleCreateContainer = async () => {
    try {
      setActionLoading({ create: true });
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
      setError('Failed to create container: ' + (err.response?.data?.message || err.message));
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
      await containerApi.startContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to start container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleStopContainer = async (containerId) => {
    try {
      setActionLoading({ [containerId]: 'stop' });
      await containerApi.stopContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to stop container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleRestartContainer = async (containerId) => {
    try {
      setActionLoading({ [containerId]: 'restart' });
      await containerApi.restartContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to restart container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const handleDeleteContainer = async (containerId) => {
    if (!window.confirm('Are you sure you want to delete this container?')) {
      return;
    }

    try {
      setActionLoading({ [containerId]: 'delete' });
      await containerApi.deleteContainer(containerId);
      await loadContainers();
    } catch (err) {
      setError('Failed to delete container: ' + (err.response?.data?.message || err.message));
    } finally {
      setActionLoading({});
    }
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'running':
        return 'success';
      case 'stopped':
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
                    <TableCell>Created At</TableCell>
                    <TableCell align="right">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {containers.map((container) => (
                    <TableRow key={container.id}>
                      <TableCell>{container.name}</TableCell>
                      <TableCell>
                        <Chip
                          label={container.status}
                          color={getStatusColor(container.status)}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>{container.state}</TableCell>
                      <TableCell>
                        {container.createdAt
                          ? new Date(container.createdAt).toLocaleString()
                          : 'N/A'}
                      </TableCell>
                      <TableCell align="right">
                        <Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
                          <Tooltip title="Start">
                            <IconButton
                              color="success"
                              size="small"
                              onClick={() => handleStartContainer(container.id)}
                              disabled={!!actionLoading[container.id]}
                            >
                              {actionLoading[container.id] === 'start' ? (
                                <CircularProgress size={20} />
                              ) : (
                                <PlayArrowIcon />
                              )}
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Stop">
                            <IconButton
                              color="error"
                              size="small"
                              onClick={() => handleStopContainer(container.id)}
                              disabled={!!actionLoading[container.id]}
                            >
                              {actionLoading[container.id] === 'stop' ? (
                                <CircularProgress size={20} />
                              ) : (
                                <StopIcon />
                              )}
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Restart">
                            <IconButton
                              color="primary"
                              size="small"
                              onClick={() => handleRestartContainer(container.id)}
                              disabled={!!actionLoading[container.id]}
                            >
                              {actionLoading[container.id] === 'restart' ? (
                                <CircularProgress size={20} />
                              ) : (
                                <RefreshIcon />
                              )}
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Delete">
                            <IconButton
                              color="error"
                              size="small"
                              onClick={() => handleDeleteContainer(container.id)}
                              disabled={!!actionLoading[container.id]}
                            >
                              {actionLoading[container.id] === 'delete' ? (
                                <CircularProgress size={20} />
                              ) : (
                                <DeleteIcon />
                              )}
                            </IconButton>
                          </Tooltip>
                        </Box>
                      </TableCell>
                    </TableRow>
                  ))}
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
              onChange={(e) => setNewContainer({ ...newContainer, name: e.target.value })}
              helperText="A unique name for your container"
            />
            <TextField
              label="Docker Image"
              fullWidth
              variant="outlined"
              required
              value={newContainer.image}
              onChange={(e) => setNewContainer({ ...newContainer, image: e.target.value })}
              helperText="e.g., nginx:latest, ubuntu:22.04, postgres:15"
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
            disabled={!newContainer.name || !newContainer.image || actionLoading.create}
          >
            {actionLoading.create ? <CircularProgress size={24} /> : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default DashboardPage;
