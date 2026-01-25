import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  CardActions,
  Button,
  CircularProgress,
  Chip,
  Stack,
  Grid,
  Paper,
  Tabs,
  Tab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  LinearProgress,
  Tooltip,
  ButtonGroup,
} from '@mui/material';
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import LockIcon from '@mui/icons-material/Lock';
import LockOpenIcon from '@mui/icons-material/LockOpen';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CloseIcon from '@mui/icons-material/Close';
import FilterListIcon from '@mui/icons-material/FilterList';
import ViewModuleIcon from '@mui/icons-material/ViewModule';
import AccountTree from '@mui/icons-material/AccountTree';
import RefreshIcon from '@mui/icons-material/Refresh';
import cytoscape from 'cytoscape';
import coseBilkent from 'cytoscape-cose-bilkent';
import { apiCall } from '../api';
import { showToast } from '../utils/toast';

cytoscape.use(coseBilkent);

function RoadmapTab() {
  const [questTree, setQuestTree] = useState([]);
  const [branches, setBranches] = useState([]);
  const [passiveBuffs, setPassiveBuffs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedBranch, setSelectedBranch] = useState('all');
  const [selectedQuest, setSelectedQuest] = useState(null);
  const [questDialogOpen, setQuestDialogOpen] = useState(false);
  const [viewMode, setViewMode] = useState('graph'); // 'graph' or 'grid'
  const cyRef = useRef(null);
  const containerRef = useRef(null);

  useEffect(() => {
    loadData();
    return () => {
      if (cyRef.current) {
        cyRef.current.destroy();
      }
    };
  }, []);

  useEffect(() => {
    if (viewMode === 'graph' && questTree.length > 0) {
      renderQuestGraph();
    }
  }, [viewMode, questTree, selectedBranch]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [treeResult, branchesResult, buffsResult] = await Promise.all([
        apiCall('/quests/tree'),
        apiCall('/branches'),
        apiCall('/quests/buffs'),
      ]);

      // Логирование для отладки
      console.log('🔍 Roadmap API Results:', {
        tree: { success: treeResult.success, dataType: typeof treeResult.data, isArray: Array.isArray(treeResult.data), length: Array.isArray(treeResult.data) ? treeResult.data.length : 'N/A', data: treeResult.data },
        branches: { success: branchesResult.success, dataType: typeof branchesResult.data, isArray: Array.isArray(branchesResult.data), length: Array.isArray(branchesResult.data) ? branchesResult.data.length : 'N/A', data: branchesResult.data },
        buffs: { success: buffsResult.success, dataType: typeof buffsResult.data, isArray: Array.isArray(buffsResult.data), length: Array.isArray(buffsResult.data) ? buffsResult.data.length : 'N/A', data: buffsResult.data },
      });

      // API возвращает данные через наш apiCall wrapper: {success: bool, data: any}
      // Gin возвращает JSON напрямую, но наш apiCall оборачивает в {success, data}
      let treeData = [];
      let branchesData = [];
      let buffsData = [];

      // Обработка treeResult
      if (treeResult.success) {
        if (Array.isArray(treeResult.data)) {
          treeData = treeResult.data;
          console.log('✅ Tree data loaded:', treeData.length, 'quests');
        } else if (treeResult.data && typeof treeResult.data === 'object') {
          // Если это объект, возможно это массив внутри или ошибка
          console.warn('⚠️ Tree data is not an array:', treeResult.data);
          // Попробуем извлечь массив из объекта
          if (treeResult.data.quests && Array.isArray(treeResult.data.quests)) {
            treeData = treeResult.data.quests;
            console.log('✅ Found quests array in object:', treeData.length);
          } else if (treeResult.data.data && Array.isArray(treeResult.data.data)) {
            treeData = treeResult.data.data;
            console.log('✅ Found data array in object:', treeData.length);
          }
        } else {
          console.warn('⚠️ Tree data is not an array or object:', typeof treeResult.data, treeResult.data);
        }
      } else {
        console.error('❌ Tree API error:', treeResult);
        const errorMsg = treeResult.data?.error || treeResult.data?.message || 'Неизвестная ошибка';
        showToast('Ошибка загрузки дерева квестов: ' + errorMsg, 'error');
      }

      // Обработка branchesResult
      if (branchesResult.success) {
        if (Array.isArray(branchesResult.data)) {
          branchesData = branchesResult.data;
        } else {
          console.warn('Branches data is not an array:', branchesResult.data);
        }
      } else {
        console.error('Branches API error:', branchesResult.data);
      }

      // Обработка buffsResult
      if (buffsResult.success) {
        if (Array.isArray(buffsResult.data)) {
          buffsData = buffsResult.data;
        } else {
          console.warn('Buffs data is not an array:', buffsResult.data);
        }
      } else {
        console.error('Buffs API error:', buffsResult.data);
      }

      setQuestTree(treeData);
      setBranches(branchesData);
      setPassiveBuffs(buffsData);

      // Детальное логирование для отладки
      console.log('📊 Processed Data:', {
        treeData: { length: treeData.length, sample: treeData[0] },
        branchesData: { length: branchesData.length, sample: branchesData[0] },
        buffsData: { length: buffsData.length, sample: buffsData[0] },
      });

      if (treeData.length === 0 && branchesData.length === 0) {
        console.error('❌ No data loaded! Check API responses above.');
        showToast('Не удалось загрузить данные. Откройте консоль (F12) для деталей.', 'warning');
      } else if (treeData.length === 0) {
        console.warn('⚠️ Tree data is empty, but branches loaded');
        showToast('Квесты не найдены. Возможно, в базе данных нет квестов.', 'info');
      }
    } catch (error) {
      showToast('Ошибка загрузки данных: ' + error.message, 'error');
      console.error('Roadmap load error:', error);
    } finally {
      setLoading(false);
    }
  };

  const filteredQuests = selectedBranch === 'all'
    ? questTree
    : questTree.filter(q => q.branches?.some(b => b.name === selectedBranch || b.name.startsWith(selectedBranch + '_')));

  const getQuestStatus = (quest) => {
    // Backend возвращает is_unlocked (JSON тег), но может быть IsUnlocked в некоторых случаях
    const isUnlocked = quest.is_unlocked !== undefined ? quest.is_unlocked : quest.IsUnlocked;
    if (isUnlocked === false) return 'locked';
    return 'available';
  };

  const getBranchColor = (branchName) => {
    const branch = branches.find(b => b.name === branchName);
    return branch?.color || '#666666';
  };

  const getBranchIcon = (branchName) => {
    const branch = branches.find(b => b.name === branchName);
    return branch?.icon || 'circle';
  };

  const handleQuestClick = (quest) => {
    if (!quest) {
      console.error('Cannot open quest dialog: quest is null');
      return;
    }
    setSelectedQuest(quest);
    setQuestDialogOpen(true);
  };

  const handlePurchaseQuest = async (questId) => {
    try {
      const result = await apiCall(`/quests/${questId}/purchase`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Квест куплен!', 'success');
        loadData();
        setQuestDialogOpen(false);
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка покупки квеста', 'error');
    }
  };

  const renderQuestGraph = () => {
    if (!containerRef.current) return;

    if (cyRef.current) {
      cyRef.current.destroy();
    }

    const filtered = selectedBranch === 'all'
      ? questTree
      : questTree.filter(q => q.branches?.some(b => b.name === selectedBranch || b.name.startsWith(selectedBranch + '_')));

    if (filtered.length === 0) return;

    const elements = [];
    const questNodes = new Map();

    // Создаем узлы для квестов
    filtered.forEach(quest => {
      const questData = quest.quest || quest;
      if (!questData || !questData.id) return;

      const isUnlocked = quest.is_unlocked !== undefined ? quest.is_unlocked : (quest.IsUnlocked !== undefined ? quest.IsUnlocked : true);
      const questNodeId = `quest_${questData.id}`;
      
      questNodes.set(questData.id, {
        nodeId: questNodeId,
        quest: quest,
        questData: questData,
        isUnlocked: isUnlocked,
      });

      const mainBranch = quest.branches?.find(b => b.level === 1);
      const branchColor = mainBranch?.color || '#666666';
      const statusColor = isUnlocked ? '#10b981' : '#ef4444';

      elements.push({
        data: {
          id: questNodeId,
          label: questData.title || 'Без названия',
          questId: questData.id,
          quest: quest,
          isUnlocked: isUnlocked,
          branchColor: branchColor,
          level: questData.quest_level || 1,
        },
        style: {
          'background-color': statusColor,
          'border-color': branchColor,
          'border-width': 3,
        },
      });
    });

    // Создаем связи между квестами (prerequisites)
    filtered.forEach(quest => {
      const questData = quest.quest || quest;
      if (!questData || !questData.id) return;

      if (quest.prerequisites && quest.prerequisites.length > 0) {
        quest.prerequisites.forEach(prereq => {
          const prereqNode = questNodes.get(prereq.prerequisite_quest_id);
          if (prereqNode) {
            elements.push({
              data: {
                id: `edge_${prereq.prerequisite_quest_id}_${questData.id}`,
                source: prereqNode.nodeId,
                target: `quest_${questData.id}`,
              },
            });
          }
        });
      }
    });

    const cy = cytoscape({
      container: containerRef.current,
      elements: elements,
      style: [
        {
          selector: 'node',
          style: {
            'label': 'data(label)',
            'color': '#ffffff',
            'font-size': '12px',
            'font-weight': '600',
            'text-valign': 'center',
            'text-halign': 'center',
            'width': 'label',
            'height': 'label',
            'padding': '12px',
            'shape': 'roundrectangle',
            'border-width': 3,
            'background-color': 'data(background-color)',
            'border-color': 'data(border-color)',
            'text-wrap': 'wrap',
            'text-max-width': '150px',
          },
        },
        {
          selector: 'edge',
          style: {
            'width': 2,
            'line-color': '#9ca3af',
            'target-arrow-color': '#9ca3af',
            'target-arrow-shape': 'triangle',
            'curve-style': 'bezier',
            'arrow-scale': 1.5,
          },
        },
        {
          selector: 'node[isUnlocked = true]',
          style: {
            'background-color': '#10b981',
          },
        },
        {
          selector: 'node[isUnlocked = false]',
          style: {
            'background-color': '#ef4444',
            'opacity': 0.6,
          },
        },
      ],
      layout: {
        name: 'cose-bilkent',
        nodeDimensionsIncludeLabels: true,
        idealEdgeLength: 100,
        nodeRepulsion: 4500,
        nestingFactor: 0.1,
        gravity: 0.25,
        numIter: 2500,
        tile: true,
        animate: true,
        animationDuration: 1000,
      },
    });

    cy.on('tap', 'node', (evt) => {
      const node = evt.target;
      const quest = node.data('quest');
      if (quest) {
        handleQuestClick(quest);
      }
    });

    cyRef.current = cy;
  };

  const mainBranches = branches.filter(b => b.level === 1);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress size={60} />
      </Box>
    );
  }

  return (
    <Box>
      <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 4 }}>
        <AccountTreeIcon sx={{ fontSize: 48, color: 'primary.main' }} />
        <Box>
          <Typography variant="h3" gutterBottom sx={{ fontWeight: 700 }}>
            Roadmap развития
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Дерево квестов и ваш путь к совершенству
          </Typography>
        </Box>
      </Stack>

      {/* Фильтр по веткам и переключатель вида */}
      <Paper elevation={2} sx={{ p: 2, mb: 3, bgcolor: 'background.paper' }}>
        <Stack direction="row" spacing={2} alignItems="center" justifyContent="space-between" sx={{ mb: 2 }}>
          <Stack direction="row" spacing={2} alignItems="center">
            <FilterListIcon />
            <Typography variant="h6">Фильтр по веткам:</Typography>
          </Stack>
          <ButtonGroup variant="outlined" size="small">
            <Button
              startIcon={<AccountTree />}
              onClick={() => setViewMode('graph')}
              variant={viewMode === 'graph' ? 'contained' : 'outlined'}
            >
              Граф
            </Button>
            <Button
              startIcon={<ViewModuleIcon />}
              onClick={() => setViewMode('grid')}
              variant={viewMode === 'grid' ? 'contained' : 'outlined'}
            >
              Сетка
            </Button>
          </ButtonGroup>
        </Stack>
        <Stack direction="row" spacing={1} flexWrap="wrap">
          <Chip
            label="Все"
            onClick={() => setSelectedBranch('all')}
            color={selectedBranch === 'all' ? 'primary' : 'default'}
            sx={{ cursor: 'pointer' }}
          />
          {mainBranches.map(branch => (
            <Chip
              key={branch.id}
              label={branch.display_name}
              onClick={() => setSelectedBranch(branch.name)}
              color={selectedBranch === branch.name ? 'primary' : 'default'}
              sx={{
                cursor: 'pointer',
                borderColor: branch.color,
                '&:hover': {
                  borderColor: branch.color,
                },
              }}
            />
          ))}
        </Stack>
      </Paper>

      {/* Пассивные баффы */}
      {passiveBuffs.length > 0 && (
        <Paper elevation={2} sx={{ p: 2, mb: 3, bgcolor: 'background.paper' }}>
          <Typography variant="h6" gutterBottom>
            Активные пассивные баффы
          </Typography>
          <Grid container spacing={2}>
            {passiveBuffs.map(buff => (
              <Grid item xs={12} sm={6} md={4} key={buff.id}>
                <Card sx={{ bgcolor: 'background.default', border: '1px solid', borderColor: 'divider' }}>
                  <CardContent>
                    <Typography variant="subtitle2" color="text.secondary">
                      {buff.buff_type === 'reward_multiplier' ? 'Множитель награды' : 'Усиление'}
                    </Typography>
                    {buff.buff_data && (
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        {(() => {
                          try {
                            // buff_data может быть строкой или уже объектом
                            const data = typeof buff.buff_data === 'string' 
                              ? JSON.parse(buff.buff_data) 
                              : buff.buff_data;
                            return data.description || data.type || 'Пассивный бафф';
                          } catch (e) {
                            console.error('Error parsing buff_data:', e, buff.buff_data);
                            return 'Пассивный бафф';
                          }
                        })()}
                      </Typography>
                    )}
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Paper>
      )}

      {/* Дерево квестов */}
      {filteredQuests.length === 0 ? (
        <Paper elevation={2} sx={{ p: 4, textAlign: 'center', bgcolor: 'background.paper' }}>
          <Typography variant="h6" color="text.secondary">
            Квесты не найдены
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Попробуйте изменить фильтр или создайте квест через AI
          </Typography>
        </Paper>
      ) : viewMode === 'graph' ? (
        <Paper elevation={2} sx={{ p: 2, bgcolor: 'background.paper' }}>
          <Box
            ref={containerRef}
            sx={{
              width: '100%',
              height: '600px',
              border: '1px solid',
              borderColor: 'divider',
              borderRadius: 2,
              bgcolor: '#f9fafb',
            }}
          />
        </Paper>
      ) : (
        <Grid container spacing={3}>
        {filteredQuests.map((quest, index) => {
          if (!quest) {
            console.warn('Invalid quest at index', index, ':', quest);
            return null;
          }
          // Проверяем структуру данных - может быть quest.quest или просто quest
          const questData = quest.quest || quest;
          if (!questData || !questData.id) {
            console.warn('Invalid quest data structure at index', index, ':', quest);
            return null;
          }
          const status = getQuestStatus(quest);
          // Backend возвращает is_unlocked (JSON тег)
          const isUnlocked = quest.is_unlocked !== undefined ? quest.is_unlocked : (quest.IsUnlocked !== undefined ? quest.IsUnlocked : true);
          const isLocked = !isUnlocked;
          const mainBranch = quest.branches?.find(b => b.level === 1);

          return (
            <Grid item xs={12} sm={6} md={4} key={questData.id}>
              <Card
                sx={{
                  bgcolor: 'background.paper',
                  border: '2px solid',
                  borderColor: isLocked ? 'error.main' : mainBranch?.color || 'divider',
                  opacity: isLocked ? 0.6 : 1,
                  cursor: 'pointer',
                  transition: 'all 0.3s ease',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: 6,
                  },
                }}
                onClick={() => handleQuestClick(quest)}
              >
                <CardContent>
                  <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
                    {isLocked ? (
                      <LockIcon color="error" />
                    ) : (
                      <LockOpenIcon color="success" />
                    )}
                    <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 600 }}>
                      {questData.title || 'Без названия'}
                    </Typography>
                  </Stack>

                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2, minHeight: '40px' }}>
                    {questData.description ? (questData.description.length > 100 ? questData.description.substring(0, 100) + '...' : questData.description) : 'Нет описания'}
                  </Typography>

                  <Stack direction="row" spacing={1} flexWrap="wrap" sx={{ mb: 2 }}>
                    {quest.branches?.slice(0, 3).map(branch => (
                      <Chip
                        key={branch.id}
                        label={branch.display_name}
                        size="small"
                        sx={{
                          bgcolor: branch.color + '20',
                          color: branch.color,
                          border: `1px solid ${branch.color}40`,
                        }}
                      />
                    ))}
                  </Stack>

                  {isLocked && quest.prerequisites && quest.prerequisites.length > 0 && (quest.required_count || quest.RequiredCount) > 0 && (
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="caption" color="text.secondary">
                        Требуется выполнить: {quest.unlocked_count || quest.UnlockedCount || 0}/{quest.required_count || quest.RequiredCount || 0} квестов
                      </Typography>
                      <LinearProgress
                        variant="determinate"
                        value={((quest.required_count || quest.RequiredCount || 0) > 0) ? (((quest.unlocked_count || quest.UnlockedCount || 0) / (quest.required_count || quest.RequiredCount || 0)) * 100) : 0}
                        sx={{ mt: 1, height: 6, borderRadius: 3 }}
                      />
                    </Box>
                  )}

                  <Stack direction="row" spacing={2} justifyContent="space-between" sx={{ mt: 2 }}>
                    <Chip
                      label={questData.rarity || 'common'}
                      size="small"
                      color={questData.rarity === 'epic' ? 'secondary' : 'default'}
                    />
                    <Typography variant="body2" color="text.secondary">
                      {questData.price || 0} монет
                    </Typography>
                  </Stack>
                </CardContent>
              </Card>
            </Grid>
          );
        })}
        </Grid>
      )}

      {/* Диалог с деталями квеста */}
      <Dialog
        open={questDialogOpen}
        onClose={() => setQuestDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          <Stack direction="row" justifyContent="space-between" alignItems="center">
            <Typography variant="h5">{selectedQuest?.quest?.title || 'Детали квеста'}</Typography>
            <IconButton onClick={() => setQuestDialogOpen(false)}>
              <CloseIcon />
            </IconButton>
          </Stack>
        </DialogTitle>
        <DialogContent dividers>
          {selectedQuest && selectedQuest.quest && (
            <Stack spacing={3}>
              <Typography variant="body1" color="text.secondary">
                {selectedQuest.quest?.description || 'Нет описания'}
              </Typography>

              <Box>
                <Typography variant="subtitle2" gutterBottom>
                  Ветки развития:
                </Typography>
                <Stack direction="row" spacing={1} flexWrap="wrap">
                  {selectedQuest.branches?.map(branch => (
                    <Chip
                      key={branch.id}
                      label={branch.display_name}
                      sx={{
                        bgcolor: branch.color + '20',
                        color: branch.color,
                        border: `1px solid ${branch.color}40`,
                      }}
                    />
                  ))}
                </Stack>
              </Box>

              {selectedQuest.prerequisites && selectedQuest.prerequisites.length > 0 && (
                <Box>
                  <Typography variant="subtitle2" gutterBottom>
                    Требования:
                  </Typography>
                  {(() => {
                    const isUnlocked = selectedQuest.is_unlocked !== undefined ? selectedQuest.is_unlocked : (selectedQuest.IsUnlocked !== undefined ? selectedQuest.IsUnlocked : true);
                    const requiredCount = selectedQuest.required_count || selectedQuest.RequiredCount || 0;
                    const unlockedCount = selectedQuest.unlocked_count || selectedQuest.UnlockedCount || 0;
                    return !isUnlocked && requiredCount > 0 && (
                      <Box sx={{ mb: 2 }}>
                        <Typography variant="body2" color="error">
                          Выполнено: {unlockedCount} из {requiredCount}
                        </Typography>
                        <LinearProgress
                          variant="determinate"
                          value={requiredCount > 0 ? (unlockedCount / requiredCount) * 100 : 0}
                          sx={{ mt: 1, height: 8, borderRadius: 4 }}
                        />
                      </Box>
                    );
                  })()}
                  <Typography variant="body2" color="text.secondary">
                    Нужно завершить {selectedQuest.required_count || selectedQuest.RequiredCount || 0} из следующих квестов:
                  </Typography>
                  <Stack spacing={1} sx={{ mt: 1 }}>
                    {selectedQuest.prerequisites.map((prereq, idx) => (
                      <Chip
                        key={idx}
                        label={`Квест #${prereq.prerequisite_quest_id}`}
                        size="small"
                        color="default"
                      />
                    ))}
                  </Stack>
                </Box>
              )}

              <Box>
                <Typography variant="subtitle2" gutterBottom>
                  Информация:
                </Typography>
                <Grid container spacing={2}>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Сложность: {selectedQuest.quest?.difficulty || 0}/10
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Редкость: {selectedQuest.quest?.rarity || 'N/A'}
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Награда XP: {selectedQuest.quest?.reward_xp || 0}
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Награда монет: {selectedQuest.quest?.reward_coin || 0}
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Цена: {selectedQuest.quest?.price || 0} монет
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Задач: {selectedQuest.quest?.tasks_count || 0}
                    </Typography>
                  </Grid>
                </Grid>
              </Box>
            </Stack>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setQuestDialogOpen(false)}>Закрыть</Button>
          {selectedQuest && selectedQuest.quest && (() => {
            const isUnlocked = selectedQuest?.is_unlocked !== undefined ? selectedQuest.is_unlocked : (selectedQuest?.IsUnlocked !== undefined ? selectedQuest.IsUnlocked : true);
            return isUnlocked && (
              <Button
                variant="contained"
                onClick={() => handlePurchaseQuest(selectedQuest.quest.id)}
              >
                Купить за {selectedQuest.quest.price || 0} монет
              </Button>
            );
          })()}
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default RoadmapTab;
