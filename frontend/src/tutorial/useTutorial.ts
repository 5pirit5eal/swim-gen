import { driver } from 'driver.js'
import 'driver.js/dist/driver.css'
import { useI18n } from 'vue-i18n'
import { useSettingsStore } from '@/stores/settings'
import { nextTick } from 'vue'

export function useTutorial() {
  const { t } = useI18n()
  const settings = useSettingsStore()

  function startHomeTutorial() {
    if (settings.tutorials.home) return

    const driverObj = driver({
      showProgress: true,
      animate: true,
      doneBtnText: t('tutorial.done'),
      nextBtnText: t('tutorial.next'),
      prevBtnText: t('tutorial.prev'),
      allowClose: true,
      onDestroyed: () => {
        settings.markTutorialSeen('home')
      },
      steps: [
        {
          element: '#tutorial-sidebar-toggle',
          popover: {
            title: t('tutorial.home.sidebar.title'),
            description: t('tutorial.home.sidebar.description'),
            side: 'bottom',
            align: 'start',
          },
        },
        {
          element: '#tutorial-profile-btn',
          popover: {
            title: t('tutorial.home.profile.title'),
            description: t('tutorial.home.profile.description'),
            side: 'bottom',
            align: 'end',
          },
        },
      ],
    })

    // Wait for DOM to be ready
    nextTick(() => {
      driverObj.drive()
    })
  }

  function startInteractionTutorial() {
    if (settings.tutorials.interaction) return

    const driverObj = driver({
      showProgress: true,
      animate: true,
      doneBtnText: t('tutorial.done'),
      nextBtnText: t('tutorial.next'),
      prevBtnText: t('tutorial.prev'),
      allowClose: true,
      onDestroyed: () => {
        settings.markTutorialSeen('interaction')
      },
      steps: [
        {
          element: '#tutorial-plan-display',
          popover: {
            title: t('tutorial.interaction.plan.title'),
            description: t('tutorial.interaction.plan.description'),
            side: 'top',
            align: 'center',
          },
        },
        {
          element: '#tutorial-tab-switcher',
          popover: {
            title: t('tutorial.interaction.tabs.title'),
            description: t('tutorial.interaction.tabs.description'),
            side: 'bottom',
            align: 'start',
          },
        },
        {
          element: '#tutorial-rate-btn',
          popover: {
            title: t('tutorial.interaction.rate.title'),
            description: t('tutorial.interaction.rate.description'),
            side: 'bottom',
            align: 'end',
          },
        },
        {
          element: '#tutorial-share-btn',
          popover: {
            title: t('tutorial.interaction.share.title'),
            description: t('tutorial.interaction.share.description'),
            side: 'top',
            align: 'end',
          },
        },
      ],
    })

    nextTick(() => {
      driverObj.drive()
    })
  }

  function startSidebarTutorial() {
    if (settings.tutorials.sidebar) return

    const driverObj = driver({
      showProgress: true,
      animate: true,
      doneBtnText: t('tutorial.done'),
      nextBtnText: t('tutorial.next'),
      prevBtnText: t('tutorial.prev'),
      allowClose: true,
      onDestroyed: () => {
        settings.markTutorialSeen('sidebar')
      },
      steps: [
        {
          element: '#tutorial-history-generated',
          popover: {
            title: t('tutorial.sidebar.history.title'),
            description: t('tutorial.sidebar.history.description'),
            side: 'bottom',
            align: 'start',
          },
        },
        {
          element: '#tutorial-new-plan-btn',
          popover: {
            title: t('tutorial.sidebar.new_plan.title'),
            description: t('tutorial.sidebar.new_plan.description'),
            side: 'bottom',
            align: 'center',
          },
        },
        {
          element: '#tutorial-upload-btn',
          popover: {
            title: t('tutorial.sidebar.upload.title'),
            description: t('tutorial.sidebar.upload.description'),
            side: 'bottom',
            align: 'center',
          },
        },
      ],
    })

    nextTick(() => {
      driverObj.drive()
    })
  }

  return {
    startHomeTutorial,
    startInteractionTutorial,
    startSidebarTutorial,
  }
}
