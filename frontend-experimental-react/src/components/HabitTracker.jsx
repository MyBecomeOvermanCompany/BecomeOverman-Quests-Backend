import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Checkbox,
  Stack,
  Chip,
  LinearProgress,
  Button,
  Alert,
  Tooltip,
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import { apiCall } from '../api';
import { showToast } from '../utils/toast';

function HabitTracker({ questId, taskId, taskTitle, onComplete }) {
  const [progress, setProgress] = useState(null);
  const [loading, setLoading] = useState(true);
  const [marking, setMarking] = useState(false);

  useEffect(() => {
    loadProgress();
  }, [questId, taskId]);

  const loadProgress = async () => {
    try {
      const result = await apiCall(`/quests/${questId}/${taskId}/habit-progress`);
      console.log('HabitTracker API response:', result);
      if (result.success && result.data) {
        console.log('HabitTracker progress data:', result.data);
        setProgress(result.data);
      } else {
        console.warn('HabitTracker: API returned no data or error:', result);
        setProgress({ has_requirement: false });
      }
    } catch (error) {
      console.error('HabitTracker: Error loading habit progress:', error);
      setProgress({ has_requirement: false });
    } finally {
      setLoading(false);
    }
  };

  const freezeMissedDay = async (dateStr) => {
    setMarking(true);
    try {
      const result = await apiCall(`/quests/${questId}/${taskId}/freeze-day`, {
        method: 'POST',
        body: JSON.stringify({ date: dateStr }),
      });

      if (result.success) {
        showToast('День заморожен! Списано 50 монет.', 'success');
        await loadProgress();
      } else {
        showToast('Ошибка: ' + (result.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка заморозки дня', 'error');
      console.error('Error freezing day:', error);
    } finally {
      setMarking(false);
    }
  };

  const markAsCompleted = async () => {
    setMarking(true);
    try {
      if (!progress || !progress.has_requirement) {
        // Если нет требований habit tracking, используем обычный endpoint для выполнения задачи
        const result = await apiCall(`/quests/${questId}/${taskId}/complete`, {
          method: 'POST',
        });

        if (result.success) {
          showToast('Задача выполнена!', 'success');
          // Вызываем callback для обновления UI
          if (onComplete) {
            onComplete();
          }
        } else {
          showToast('Ошибка: ' + (result.data?.error || result.error || 'Неизвестная ошибка'), 'error');
        }
        return;
      }

      // Если есть требования habit tracking, используем специальный endpoint
      const now = new Date();
      const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`;
      
      const result = await apiCall(`/quests/${questId}/${taskId}/habit-complete`, {
        method: 'POST',
        body: JSON.stringify({ completion_time: timeStr }),
      });

      if (result.success) {
        showToast('Задача отмечена как выполненная!', 'success');
        await loadProgress();
        
        // Если задача полностью выполнена, вызываем callback
        if (progress.is_completed && onComplete) {
          onComplete();
        }
      } else {
        showToast('Ошибка: ' + (result.data?.error || result.error || 'Неизвестная ошибка'), 'error');
      }
    } catch (error) {
      showToast('Ошибка выполнения задачи', 'error');
      console.error('Error completing task:', error);
    } finally {
      setMarking(false);
    }
  };

  if (loading) {
    return <Typography variant="body2">Загрузка...</Typography>;
  }

  if (!progress || !progress.has_requirement) {
    // Если нет требований habit tracking, показываем обычную кнопку
    console.log('HabitTracker: No requirement, showing simple button. Progress:', progress);
    console.log('HabitTracker: questId=', questId, 'taskId=', taskId);
    return (
      <Box sx={{ mt: 1 }}>
        <Button
          variant="contained"
          onClick={markAsCompleted}
          disabled={marking}
          fullWidth
          sx={{ 
            bgcolor: 'primary.main',
            color: 'white',
            '&:hover': {
              bgcolor: 'primary.dark',
            }
          }}
        >
          {marking ? 'Отмечаю...' : 'Выполнить задачу'}
        </Button>
      </Box>
    );
  }

  const consecutive_days = progress.consecutive_days || 0;
  const required_days = progress.required_days || 1;
  const daytime_required = progress.daytime_required;
  const completions = progress.completions || [];
  const is_completed = progress.is_completed || false;
  const missed_days = progress.missed_days || [];
  const has_missed_days = progress.has_missed_days || false;
  
  const progressPercent = required_days > 0 ? Math.min(100, (consecutive_days / required_days) * 100) : 0;
  
  console.log('HabitTracker: Rendering with progress:', {
    consecutive_days,
    required_days,
    daytime_required,
    completionsCount: completions?.length || 0,
    is_completed,
    has_requirement: progress.has_requirement,
    missed_days,
    has_missed_days
  });

  const getDaytimeLabel = (daytime) => {
    switch (daytime) {
      case 'morning':
        return 'Утро (6:00-12:00)';
      case 'afternoon':
        return 'День (12:00-18:00)';
      case 'evening':
        return 'Вечер (18:00-6:00)';
      default:
        return 'Любое время';
    }
  };

  // Получаем последние N дней для отображения чекбоксов
  const getDaysForDisplay = () => {
    const days = [];
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    
    if (!required_days || required_days <= 0) {
      console.warn('HabitTracker: required_days is invalid:', required_days);
      return [];
    }
    
    for (let i = required_days - 1; i >= 0; i--) {
      const date = new Date(today);
      date.setDate(date.getDate() - i);
      date.setHours(0, 0, 0, 0);
      const dateStr = date.toISOString().split('T')[0];
      
      const completion = completions?.find(c => {
        if (!c || !c.completion_date) return false;
        const compDate = new Date(c.completion_date);
        compDate.setHours(0, 0, 0, 0);
        const compDateStr = compDate.toISOString().split('T')[0];
        return compDateStr === dateStr;
      });
      
      days.push({
        date: dateStr,
        dateObj: date,
        isCompleted: !!completion,
        isToday: i === 0,
      });
    }
    
    console.log('HabitTracker: Generated days array:', days);
    return days;
  };

  const days = getDaysForDisplay();
  
  if (!days || days.length === 0) {
    console.warn('HabitTracker: No days generated!', { required_days, completions });
  }

  return (
    <Box sx={{ mt: 2 }}>
      <Stack spacing={2}>
        <Box>
          <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 1 }}>
            <Typography variant="subtitle2" fontWeight={600}>
              Прогресс выполнения:
            </Typography>
            <Chip
              label={`${consecutive_days} / ${required_days} дней`}
              size="small"
              color={is_completed ? 'success' : 'default'}
            />
          </Stack>
          <LinearProgress
            variant="determinate"
            value={progressPercent}
            sx={{ height: 8, borderRadius: 4 }}
            color={is_completed ? 'success' : 'primary'}
          />
        </Box>

        {daytime_required && daytime_required !== 'any' && (
          <Alert severity="info" icon={<AccessTimeIcon />} sx={{ py: 0.5 }}>
            <Typography variant="caption">
              Требуется выполнение: {getDaytimeLabel(daytime_required)}
            </Typography>
          </Alert>
        )}

        {has_missed_days && missed_days.length > 0 && (
          <Alert 
            severity="warning" 
            icon={<AttachMoneyIcon />}
            sx={{ py: 0.5 }}
            action={
              <Button
                size="small"
                color="inherit"
                onClick={() => {
                  // Замораживаем первый пропущенный день
                  if (missed_days.length > 0) {
                    freezeMissedDay(missed_days[0]);
                  }
                }}
                disabled={marking}
              >
                Заморозить (50 монет)
              </Button>
            }
          >
            <Typography variant="caption">
              Пропущено дней: {missed_days.length}. Можно заморозить за 50 монет за день.
            </Typography>
          </Alert>
        )}

        <Box>
          <Typography variant="subtitle2" gutterBottom>
            Выполнение по дням ({required_days} {required_days === 1 ? 'день' : required_days < 5 ? 'дня' : 'дней'}):
          </Typography>
          <Stack 
            direction="row" 
            spacing={1} 
            flexWrap="wrap" 
            sx={{ 
              mt: 1,
              gap: 1,
            }}
          >
            {days && days.length > 0 ? days.map((day, idx) => {
              const today = new Date();
              today.setHours(0, 0, 0, 0);
              const dayStart = new Date(day.dateObj);
              dayStart.setHours(0, 0, 0, 0);
              const dayEnd = new Date(day.dateObj);
              dayEnd.setHours(23, 59, 59, 999);
              
              const isPast = dayStart < today;
              const isFuture = dayStart > today;
              const canCheck = day.isToday && !day.isCompleted && !isFuture;
              const isMissed = has_missed_days && missed_days.includes(day.date) && !day.isCompleted && isPast;

              return (
                <Tooltip
                  key={idx}
                  title={day.dateObj.toLocaleDateString('ru-RU', {
                    weekday: 'long',
                    day: 'numeric',
                    month: 'long',
                  }) + (day.isCompleted ? ' ✓ Выполнено' : canCheck ? ' - Нажмите для отметки' : isMissed ? ' - Пропущено (50 монет для заморозки)' : isFuture ? ' - Будущее' : ' - Прошло')}
                >
                  <Box
                    sx={{
                      display: 'flex',
                      flexDirection: 'column',
                      alignItems: 'center',
                      justifyContent: 'center',
                      minWidth: 70,
                      width: 70,
                      p: 1.5,
                      borderRadius: 2,
                      border: '2px solid',
                      borderColor: day.isCompleted 
                        ? 'success.main' 
                        : isMissed
                          ? 'error.main'
                        : day.isToday 
                          ? 'primary.main' 
                          : 'divider',
                      bgcolor: day.isCompleted 
                        ? 'success.light' 
                        : isMissed
                          ? 'error.light'
                        : day.isToday 
                          ? 'primary.light' 
                          : 'background.paper',
                      cursor: (canCheck || isMissed) ? 'pointer' : 'default',
                      opacity: isFuture ? 0.4 : 1,
                      transition: 'all 0.2s',
                      '&:hover': (canCheck || isMissed) ? {
                        transform: 'scale(1.1)',
                        boxShadow: 3,
                        borderColor: isMissed ? 'error.dark' : 'primary.dark',
                      } : {},
                    }}
                    onClick={canCheck ? markAsCompleted : isMissed ? () => freezeMissedDay(day.date) : undefined}
                  >
                    <Checkbox
                      checked={day.isCompleted}
                      disabled={!canCheck && !day.isCompleted}
                      size="medium"
                      sx={{
                        p: 0,
                        color: day.isCompleted ? 'success.main' : day.isToday ? 'primary.main' : 'text.secondary',
                        '&.Mui-checked': {
                          color: 'success.main',
                        },
                        '&.Mui-disabled': {
                          color: 'text.disabled',
                        },
                      }}
                    />
                    <Typography 
                      variant="caption" 
                      color={day.isCompleted ? 'success.main' : day.isToday ? 'primary.main' : 'text.secondary'}
                      sx={{ 
                        fontSize: '0.65rem',
                        fontWeight: day.isToday ? 600 : 400,
                        mt: 0.5,
                        textAlign: 'center',
                      }}
                    >
                      {day.dateObj.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })}
                    </Typography>
                    {day.isToday && (
                      <Typography 
                        variant="caption" 
                        color="primary.main" 
                        sx={{ 
                          fontSize: '0.6rem', 
                          mt: 0.25,
                          fontWeight: 600,
                        }}
                      >
                        Сегодня
                      </Typography>
                    )}
                    {isMissed && (
                      <Typography 
                        variant="caption" 
                        color="error.main" 
                        sx={{ 
                          fontSize: '0.6rem', 
                          mt: 0.25,
                          fontWeight: 600,
                        }}
                      >
                        Пропущено
                      </Typography>
                    )}
                  </Box>
                </Tooltip>
              );
            }) : (
              <Typography variant="body2" color="text.secondary">
                Нет дней для отображения
              </Typography>
            )}
          </Stack>
        </Box>

        {!is_completed && (
          <Button
            variant="contained"
            onClick={markAsCompleted}
            disabled={marking}
            fullWidth
            startIcon={<CheckCircleIcon />}
            sx={{ mt: 1 }}
          >
            {marking ? 'Отмечаю...' : 'Отметить выполнение сегодня'}
          </Button>
        )}

        {is_completed && (
          <Alert severity="success">
            Задача выполнена! Выполнено {consecutive_days} дней подряд.
          </Alert>
        )}
      </Stack>
    </Box>
  );
}

export default HabitTracker;
