import React, { useState, useEffect, useRef } from 'react';
import { 
  Box, 
  Typography, 
  Button, 
  Card, 
  CardContent, 
  Grid, 
  CircularProgress, 
  Chip,
  Stack,
  Paper,
  ButtonGroup,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
} from '@mui/material';
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import RefreshIcon from '@mui/icons-material/Refresh';
import ViewModuleIcon from '@mui/icons-material/ViewModule';
import RadioButtonUncheckedIcon from '@mui/icons-material/RadioButtonUnchecked';
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import CloseIcon from '@mui/icons-material/Close';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import cytoscape from 'cytoscape';
import coseBilkent from 'cytoscape-cose-bilkent';
import { apiCall } from '../api';
import { showToast } from '../utils/toast';

cytoscape.use(coseBilkent);

function QuestMapTab() {
  const [loading, setLoading] = useState(true);
  const cyRef = useRef(null);
  const containerRef = useRef(null);
  const [selectedNode, setSelectedNode] = useState(null);
  const [nodeDialogOpen, setNodeDialogOpen] = useState(false);

  useEffect(() => {
    loadQuestMap();

    return () => {
      if (cyRef.current) {
        cyRef.current.destroy();
      }
    };
  }, []);

  const loadQuestMap = async () => {
    setLoading(true);
    try {
      const result = await apiCall('/quests/my-quests-with-details');

      if (result.success && result.data && Array.isArray(result.data)) {
        renderQuestMap(result.data);
      } else {
        showToast('Ошибка загрузки карты квестов', 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки карты квестов', 'error');
    } finally {
      setLoading(false);
    }
  };

  const renderQuestMap = (quests) => {
    if (!containerRef.current) return;

    if (cyRef.current) {
      cyRef.current.destroy();
    }

    const elements = [];

    quests.forEach(quest => {
      const questNodeId = `quest_${quest.id}`;

      elements.push({
        data: {
          id: questNodeId,
          label: quest.title || 'Без названия',
          type: 'quest',
          questId: quest.id,
          status: quest.status || 'available',
          description: quest.description || '',
        },
      });

      if (quest.tasks && Array.isArray(quest.tasks)) {
        quest.tasks.forEach((task) => {
          const taskNodeId = `task_${quest.id}_${task.id}`;

          elements.push({
            data: {
              id: taskNodeId,
              label: task.title || 'Задача',
              type: 'task',
              questId: quest.id,
              taskId: task.id,
              status: task.status || task.completed ? 'completed' : 'active',
              description: task.description || '',
            },
          });

          elements.push({
            data: {
              id: `edge_${quest.id}_${task.id}`,
              source: questNodeId,
              target: taskNodeId,
            },
          });
        });
      }
    });

    const cy = cytoscape({
      container: containerRef.current,
      elements: elements,
      style: [
        {
          selector: 'node[type="quest"]',
          style: {
            'background-color': function(ele) {
              const status = ele.data('status');
              return status === 'completed' ? '#10b981' :
                     status === 'started' ? '#000000' :
                     '#6b7280';
            },
            'label': 'data(label)',
            'color': '#ffffff',
            'font-size': '14px',
            'font-weight': '700',
            'text-valign': 'center',
            'text-halign': 'center',
            'width': 'label',
            'height': 'label',
            'padding': '16px',
            'shape': 'roundrectangle',
            'border-width': 2,
            'border-color': '#ffffff',
          },
        },
        {
          selector: 'node[type="task"]',
          style: {
            'background-color': function(ele) {
              const status = ele.data('status');
              return status === 'completed' ? '#10b981' :
                     status === 'active' ? '#f59e0b' :
                     '#ef4444';
            },
            'label': 'data(label)',
            'color': '#ffffff',
            'font-size': '12px',
            'font-weight': '600',
            'text-valign': 'center',
            'text-halign': 'center',
            'width': 'label',
            'height': 'label',
            'padding': '12px',
            'shape': 'ellipse',
            'border-width': 2,
            'border-color': '#ffffff',
          },
        },
        {
          selector: 'edge',
          style: {
            'width': 2,
            'line-color': '#d1d5db',
            'target-arrow-color': '#d1d5db',
            'target-arrow-shape': 'triangle',
            'curve-style': 'bezier',
          },
        },
      ],
      layout: {
        name: 'cose-bilkent',
        idealEdgeLength: 100,
        nodeOverlap: 20,
        refresh: 20,
        fit: true,
        padding: 50,
        randomize: true,
        componentSpacing: 100,
        nodeRepulsion: 400000,
        edgeElasticity: 100,
        nestingFactor: 5,
        gravity: 50,
        numIter: 1000,
      },
    });

    cy.on('tap', 'node', (evt) => {
      const node = evt.target;
      const nodeData = node.data();
      setSelectedNode(nodeData);
      setNodeDialogOpen(true);
    });

    cy.nodes().forEach(node => {
      node.grabbable(true);
    });

    cyRef.current = cy;
  };

  const applyLayout = (layoutName) => {
    if (!cyRef.current) return;

    const layouts = {
      'tree': {
        name: 'breadthfirst',
        directed: true,
        padding: 50,
        spacingFactor: 1.5,
      },
      'grid': {
        name: 'grid',
        rows: 3,
        padding: 50,
      },
      'circle': {
        name: 'circle',
        padding: 50,
      },
      'cose': {
        name: 'cose-bilkent',
        idealEdgeLength: 100,
        nodeOverlap: 20,
        refresh: 20,
        fit: true,
        padding: 50,
      },
    };

    cyRef.current.layout(layouts[layoutName] || layouts.cose).run();
  };

  const completeTaskFromMap = async () => {
    if (!selectedNode || !selectedNode.questId || !selectedNode.taskId) return;

    try {
      const result = await apiCall(`/quests/${selectedNode.questId}/${selectedNode.taskId}/complete`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Задача выполнена!', 'success');
        setNodeDialogOpen(false);
        loadQuestMap();
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка выполнения задачи', 'error');
    }
  };

  return (
    <Box>
      <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 4 }}>
        <AccountTreeIcon sx={{ fontSize: 48, color: 'primary.main' }} />
        <Box>
          <Typography variant="h3" gutterBottom sx={{ fontWeight: 700 }}>
            Карта ваших квестов
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Интерактивная визуализация всех квестов и задач
          </Typography>
        </Box>
      </Stack>

      <Stack direction="row" spacing={2} sx={{ mb: 3, flexWrap: 'wrap' }}>
        <Button 
          variant="contained" 
          onClick={loadQuestMap}
          startIcon={<RefreshIcon />}
        >
          Загрузить карту
        </Button>
        <ButtonGroup variant="outlined" size="small">
          <Button onClick={() => applyLayout('tree')} startIcon={<AccountTreeIcon />}>
            Дерево
          </Button>
          <Button onClick={() => applyLayout('grid')} startIcon={<ViewModuleIcon />}>
            Сетка
          </Button>
          <Button onClick={() => applyLayout('circle')} startIcon={<RadioButtonUncheckedIcon />}>
            Круг
          </Button>
          <Button onClick={() => applyLayout('cose')} startIcon={<AutoAwesomeIcon />}>
            Авто
          </Button>
        </ButtonGroup>
      </Stack>

      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <CircularProgress size={60} />
        </Box>
      ) : (
        <Card sx={{ p: 2 }}>
          <Box
            ref={containerRef}
            sx={{
              width: '100%',
              height: '600px',
              borderRadius: 2,
              border: '1px solid',
              borderColor: 'divider',
              bgcolor: '#fafafa',
            }}
          />
        </Card>
      )}

      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Легенда
          </Typography>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6} md={3}>
              <Stack direction="row" spacing={1} alignItems="center">
                <Box sx={{ width: 24, height: 24, bgcolor: '#000000', borderRadius: 1 }} />
                <Typography variant="body2">Квест (активный)</Typography>
              </Stack>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Stack direction="row" spacing={1} alignItems="center">
                <Box sx={{ width: 24, height: 24, bgcolor: '#10b981', borderRadius: '50%' }} />
                <Typography variant="body2">Задача (выполнена)</Typography>
              </Stack>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Stack direction="row" spacing={1} alignItems="center">
                <Box sx={{ width: 24, height: 24, bgcolor: '#f59e0b', borderRadius: '50%' }} />
                <Typography variant="body2">Задача (активная)</Typography>
              </Stack>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Stack direction="row" spacing={1} alignItems="center">
                <Box sx={{ width: 24, height: 24, bgcolor: '#ef4444', borderRadius: '50%' }} />
                <Typography variant="body2">Задача (просрочена)</Typography>
              </Stack>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      <Dialog
        open={nodeDialogOpen}
        onClose={() => setNodeDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6" fontWeight={600}>
            {selectedNode?.label || 'Узел'}
          </Typography>
          <IconButton onClick={() => setNodeDialogOpen(false)}>
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          {selectedNode && (
            <Stack spacing={2}>
              {selectedNode.description && (
                <Typography variant="body1">
                  {selectedNode.description}
                </Typography>
              )}
              {selectedNode.status === 'completed' && (
                <Chip 
                  icon={<CheckCircleIcon />}
                  label="Выполнено" 
                  color="success"
                />
              )}
            </Stack>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setNodeDialogOpen(false)}>
            Закрыть
          </Button>
          {selectedNode && selectedNode.type === 'task' && selectedNode.status !== 'completed' && (
            <Button 
              variant="contained" 
              onClick={completeTaskFromMap}
            >
              Выполнить задачу
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default QuestMapTab;
