function resumeAfterAutoPostConfirmation(action) {
  const normalizedAction = String(action || 'save')
  const confirmed = {
    skipAutoPostConfirmation: true,
    skipSensitiveSettingsConfirmation: true
  }

  if (normalizedAction.startsWith('leave:')) {
    return {
      ...confirmed,
      target: 'leave',
      leaveAction: normalizedAction.slice('leave:'.length) || 'home'
    }
  }

  return {
    ...confirmed,
    target: 'save',
    leaveAction: ''
  }
}

export {
  resumeAfterAutoPostConfirmation
}
