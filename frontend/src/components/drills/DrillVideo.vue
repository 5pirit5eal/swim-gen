<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

declare global {
  interface Window {
    onYouTubeIframeAPIReady?: () => void
  }
}

const props = defineProps<{
  videoId: string
  start: number
  end: number
}>()

const playerContainer = ref<HTMLElement | null>(null)
let player: YT.Player | null = null
let loopInterval: number | null = null

const initPlayer = () => {
  if (!playerContainer.value) return

  player = new YT.Player(playerContainer.value, {
    height: '100%',
    width: '100%',
    videoId: props.videoId,
    playerVars: {
      autoplay: 1,
      mute: 1,
      controls: 1,
      start: props.start,
      end: props.end,
    },
    events: {
      onReady: (event: YT.PlayerEvent) => {
        event.target.playVideo()
        startLoopCheck()
      },
      onStateChange: (event: YT.OnStateChangeEvent) => {
        // If video ends naturally, jump back to start
        if (event.data === YT.PlayerState.ENDED) {
          player?.seekTo(props.start, true)
          player?.playVideo()
        }
      },
    },
  })
}

const startLoopCheck = () => {
  // Safety net: Check current time every 500ms to force loop
  loopInterval = window.setInterval(() => {
    if (player && player.getCurrentTime() >= props.end) {
      player.seekTo(props.start, true)
      player.playVideo()
    }
  }, 500)
}

onMounted(() => {
  // Check if YouTube API is already fully loaded
  if (window.YT && window.YT.Player) {
    initPlayer()
  } else if (!window.YT) {
    // Load YouTube API if not already present
    const tag = document.createElement('script')
    tag.src = 'https://www.youtube.com/iframe_api'
    const firstScriptTag = document.getElementsByTagName('script')[0]
    if (firstScriptTag && firstScriptTag.parentNode) {
      firstScriptTag.parentNode.insertBefore(tag, firstScriptTag)
    } else {
      document.head.appendChild(tag)
    }

    // Chain callbacks to support multiple components
    const existingCallback = window.onYouTubeIframeAPIReady
    window.onYouTubeIframeAPIReady = () => {
      existingCallback?.()
      initPlayer()
    }
  } else {
    // API script is loading but not ready yet - chain our callback
    const existingCallback = window.onYouTubeIframeAPIReady
    window.onYouTubeIframeAPIReady = () => {
      existingCallback?.()
      initPlayer()
    }
  }
})

onUnmounted(() => {
  if (loopInterval) clearInterval(loopInterval)
  if (player) player.destroy()
})
</script>

<template>
  <div class="video-wrapper">
    <div ref="playerContainer"></div>
  </div>
</template>

<style scoped>
.video-wrapper {
  position: relative;
  width: 100%;
  padding-bottom: 56.25%;
  /* 16:9 Aspect Ratio */
  height: 0;
  margin: 0 auto;
  overflow: hidden;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  background: var(--color-background-soft);
}

.video-wrapper :deep(iframe),
.video-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
</style>
