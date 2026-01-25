import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  TextField,
  InputAdornment,
  CircularProgress,
  Checkbox,
  FormControlLabel,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Chip,
  IconButton,
  Stack,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import { apiCall } from '../api';
import QuestCard from './QuestCard';
import MasonryGrid from './MasonryGrid';
import HabitTracker from './HabitTracker';
import { showToast } from '../utils/toast';

function MyQuestsTab() {
  const [activeQuests, setActiveQuests] = useState([]);
  const [completedQuests, setCompletedQuests] = useState([]);
  const [selectedQuest, setSelectedQuest] = useState(null);
  const [selectedQuestTasks, setSelectedQuestTasks] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);
  const [taskDialogOpen, setTaskDialogOpen] = useState(false);
  const [nextQuestLevel, setNextQuestLevel] = useState(null);

  useEffect(() => {
    loadMyQuests();
  }, []);

  const loadMyQuests = async () => {
    setLoading(true);
    try {
      const active = await apiCall('/quests/active');
      const completed = await apiCall('/quests/completed');

      if (active.success && active.data && Array.isArray(active.data)) {
        setActiveQuests(active.data);
      }

      if (completed.success && completed.data && Array.isArray(completed.data)) {
        setCompletedQuests(completed.data);
      }
    } catch (error) {
      showToast('Ошибка загрузки квестов', 'error');
    } finally {
      setLoading(false);
    }
  };

  const filterQuests = (quests, searchTerm) => {
    if (!searchTerm) return quests;
    if (!Array.isArray(quests)) return [];
    const term = searchTerm.toLowerCase();
    return quests.filter(quest => 
      (quest.title && quest.title.toLowerCase().includes(term)) ||
      (quest.description && quest.description.toLowerCase().includes(term))
    );
  };

  const startQuest = async (quest) => {
    if (!quest || !quest.id) return;
    
    try {
      const result = await apiCall(`/quests/${quest.id}/start`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Квест начат!', 'success');
        await loadMyQuests();
        // Автоматически показываем задачи после старта
        await showQuestTasks(quest);
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка начала квеста', 'error');
    }
  };

  const completeQuest = async (quest) => {
    if (!quest || !quest.id) return;
    
    try {
      const result = await apiCall(`/quests/${quest.id}/complete`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Квест завершен!', 'success');
        // Проверяем, есть ли следующий уровень
        await checkNextQuestLevel(quest.id);
        loadMyQuests();
        if (selectedQuest?.id === quest.id) {
          setSelectedQuest(null);
          setSelectedQuestTasks([]);
          setTaskDialogOpen(false);
        }
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка завершения квеста', 'error');
    }
  };

  const continueQuest = async (quest) => {
    if (!quest || !quest.id) return;
    
    try {
      const result = await apiCall(`/quests/${quest.id}/continue`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Следующий уровень квеста начат!', 'success');
        loadMyQuests();
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка продолжения квеста', 'error');
    }
  };

  const checkNextQuestLevel = async (questId) => {
    try {
      const result = await apiCall(`/quests/${questId}/details`);
      if (result.success && result.data) {
        const quest = result.data;
        if (quest.max_level && quest.quest_level && quest.quest_level < quest.max_level) {
          setNextQuestLevel(quest);
        } else {
          setNextQuestLevel(null);
        }
      }
    } catch (error) {
      console.error('Error checking next quest level:', error);
      setNextQuestLevel(null);
    }
  };

  const showQuestTasks = async (quest) => {
    if (!quest || !quest.id) return;
    
    try {
      const result = await apiCall(`/quests/${quest.id}/details`);
      if (result.success) {
        setSelectedQuest(quest);
        setSelectedQuestTasks(result.data.tasks || []);
        setTaskDialogOpen(true);
        // Проверяем следующий уровень, если квест завершен
        if (quest.status === 'completed') {
          await checkNextQuestLevel(quest.id);
        }
      } else {
        showToast('Ошибка загрузки деталей', 'error');
      }
    } catch (error) {
      showToast('Ошибка загрузки деталей', 'error');
    }
  };

  const showQuestDetails = async (quest) => {
    // Используем ту же функцию, что и для задач
    await showQuestTasks(quest);
  };

  const completeTask = async (taskId) => {
    if (!selectedQuest || !selectedQuest.id || !taskId) return;

    try {
      const result = await apiCall(`/quests/${selectedQuest.id}/${taskId}/complete`, {
        method: 'POST',
      });

      if (result.success) {
        showToast('Задача выполнена!', 'success');
        setSelectedQuestTasks(prev => prev.map(task => 
          task.id === taskId ? { ...task, completed: true, status: 'completed' } : task
        ));
        loadMyQuests();
      } else {
        showToast('Ошибка: ' + (result.data?.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка выполнения задачи', 'error');
    }
  };

  const handleQuestAction = (quest) => {
    if (!quest) return;
    
    if (quest.status === 'started') {
      showQuestTasks(quest);
    } else {
      startQuest(quest);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress size={60} />
      </Box>
    );
  }

  const filteredActive = filterQuests(activeQuests, searchTerm);
  const filteredCompleted = filterQuests(completedQuests, searchTerm);

  return (
    <Box sx={{ pb: 4 }}>
      <Typography variant="h3" gutterBottom sx={{ mb: 4, fontWeight: 700 }}>
        Мои квесты
      </Typography>

      <TextField
        fullWidth
        placeholder="Поиск по моим квестам..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        sx={{ mb: 4, maxWidth: 600 }}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon />
            </InputAdornment>
          ),
        }}
      />

      {filteredActive.length > 0 && (
        <Box sx={{ mb: 6 }}>
          <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
            Активные квесты
          </Typography>
          <MasonryGrid>
            {filteredActive.map(quest => (
              <QuestCard
                key={quest.id}
                quest={quest}
                onAction={handleQuestAction}
                onDetails={showQuestDetails}
                actionLabel={quest.status === 'started' ? 'Задачи' : 'Начать'}
                showProgress={true}
              />
            ))}
          </MasonryGrid>
        </Box>
      )}

      {filteredActive.length === 0 && !loading && (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            Нет активных квестов
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Начните новый квест из раздела "Квесты"!
          </Typography>
        </Box>
      )}

      {filteredCompleted.length > 0 && (
        <Box>
          <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
            Завершенные квесты
          </Typography>
          <MasonryGrid>
            {filteredCompleted.map(quest => (
              <Box key={quest.id}>
                <QuestCard
                  quest={quest}
                  completed={true}
                  onAction={() => showQuestTasks(quest)}
                  onDetails={showQuestDetails}
                  actionLabel="Просмотреть"
                />
                {quest.max_level && quest.quest_level && quest.quest_level < quest.max_level && (
                  <Box sx={{ mt: 1, textAlign: 'center' }}>
                    <Button
                      variant="outlined"
                      color="primary"
                      onClick={() => continueQuest(quest)}
                      size="small"
                    >
                      Продолжить уровень {quest.quest_level + 1}
                    </Button>
                  </Box>
                )}
              </Box>
            ))}
          </MasonryGrid>
        </Box>
      )}

      <Dialog
        open={taskDialogOpen}
        onClose={() => setTaskDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h5" fontWeight={600}>
            {selectedQuest?.title || 'Задачи квеста'}
          </Typography>
          <IconButton onClick={() => setTaskDialogOpen(false)}>
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent>
          {selectedQuestTasks.length === 0 ? (
            <Typography color="text.secondary">Задачи не найдены</Typography>
          ) : (
            <Stack spacing={2} sx={{ mt: 1 }}>
              {selectedQuestTasks.map((task, idx) => {
                const isCompleted = task.completed || task.status === 'completed';
                return (
                  <Box
                    key={task.id || idx}
                    sx={{
                      p: 2.5,
                      borderRadius: 2,
                      bgcolor: isCompleted ? '#f0fdf4' : '#fafafa',
                      border: `1px solid ${isCompleted ? '#86efac' : '#e5e7eb'}`,
                    }}
                  >
                    <FormControlLabel
                      control={
                        <Checkbox
                          checked={isCompleted}
                          onChange={() => !isCompleted && completeTask(task.id)}
                          disabled={isCompleted}
                          icon={<PlayArrowIcon />}
                          checkedIcon={<CheckCircleIcon />}
                          color="success"
                        />
                      }
                      label={
                        <Box>
                          <Typography variant="subtitle1" fontWeight={600}>
                            {task.title || `Задача ${idx + 1}`}
                          </Typography>
                          {task.description && (
                            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                              {task.description}
                            </Typography>
                          )}
                          {(task.base_xp_reward || task.base_coin_reward) && (
                            <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
                              {task.base_xp_reward && (
                                <Chip label={`${task.base_xp_reward} XP`} size="small" />
                              )}
                              {task.base_coin_reward && (
                                <Chip label={`${task.base_coin_reward} монет`} size="small" />
                              )}
                            </Stack>
                          )}
                        </Box>
                      }
                    />
                    {!isCompleted && selectedQuest && (
                      <HabitTracker
                        questId={selectedQuest.id}
                        taskId={task.id}
                        taskTitle={task.title}
                        onComplete={() => {
                          loadMyQuests();
                          showQuestTasks(selectedQuest);
                        }}
                      />
                    )}
                  </Box>
                );
              })}
            </Stack>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setTaskDialogOpen(false)}>
            Закрыть
          </Button>
          {selectedQuest && selectedQuest.status === 'completed' && nextQuestLevel && (
            <Button
              variant="contained"
              color="primary"
              onClick={() => continueQuest(selectedQuest)}
            >
              Продолжить уровень {nextQuestLevel.quest_level + 1}
            </Button>
          )}
          {selectedQuest && selectedQuest.status !== 'completed' && (
            <Button 
              variant="contained" 
              color="warning"
              onClick={() => completeQuest(selectedQuest)}
            >
              Завершить квест
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default MyQuestsTab;
