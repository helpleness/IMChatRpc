<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { PhoneOff, Mic, MicOff, Video, VideoOff } from 'lucide-vue-next';

// 定义 props 来接收通话类型和联系人信息
const props = defineProps({
  type: {
    type: String,
    validator: (value) => ['voice', 'video'].includes(value),
    required: true,
  },
  contact: {
    type: Object,
    required: true,
  },
});

const emit = defineEmits(['close']);

// 通话状态
const isMuted = ref(false);
const isCameraOn = ref(props.type === 'video'); // 视频通话默认开启摄像头
const callDuration = ref('00:00');
let timerInterval = null;

const toggleMute = () => {
  isMuted.value = !isMuted.value;
  console.log(`Mock Action: Mic is now ${isMuted.value ? 'OFF' : 'ON'}`);
};

const toggleCamera = () => {
  if (props.type === 'video') {
    isCameraOn.value = !isCameraOn.value;
    console.log(`Mock Action: Camera is now ${isCameraOn.value ? 'ON' : 'OFF'}`);
  }
};

const endCall = () => {
  console.log('Mock Action: Ending call.');
  emit('close');
};

// 在组件挂载时启动计时器
onMounted(() => {
  let seconds = 0;
  timerInterval = setInterval(() => {
    seconds++;
    const mins = Math.floor(seconds / 60).toString().padStart(2, '0');
    const secs = (seconds % 60).toString().padStart(2, '0');
    callDuration.value = `${mins}:${secs}`;
  }, 1000);
});

// 在组件卸载时清除计时器，防止内存泄漏
onUnmounted(() => {
  if (timerInterval) {
    clearInterval(timerInterval);
  }
});
</script>

<template>
  <div class="modal-backdrop">
    <div class="call-modal-content">
      <img :src="contact.avatar" :alt="contact.name" class="contact-avatar-large" />
      <h3 class="contact-name">{{ contact.name }}</h3>
      <p class="call-status">
        {{ type === 'video' ? (isCameraOn ? 'Video calling...' : 'Video call, camera off') : 'Voice calling...' }}
      </p>
      <p class="call-duration">{{ callDuration }}</p>

      <div v-if="type === 'video'" class="video-placeholder" :class="{ 'camera-off': !isCameraOn }">
        <VideoOff v-if="!isCameraOn" :size="48" color="#fff" />
        <p v-if="!isCameraOn">Camera is off</p>
      </div>

      <div class="call-controls">
        <button class="control-btn" @click="toggleMute" :class="{ 'toggled-off': isMuted }">
          <MicOff v-if="isMuted" :size="24" />
          <Mic v-else :size="24" />
        </button>
        <button v-if="type === 'video'" class="control-btn" @click="toggleCamera" :class="{ 'toggled-off': !isCameraOn }">
          <VideoOff v-if="!isCameraOn" :size="24" />
          <Video v-else :size="24" />
        </button>
        <button class="control-btn end-call-btn" @click="endCall">
          <PhoneOff :size="24" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(30, 30, 30, 0.9);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  color: white;
}

.call-modal-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.contact-avatar-large {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 3px solid white;
  box-shadow: 0 0 20px rgba(0,0,0,0.5);
}

.contact-name {
  font-size: 2rem;
  font-weight: 500;
  margin: 0;
}

.call-status, .call-duration {
  font-size: 1rem;
  margin: 0;
  opacity: 0.8;
}

.call-duration {
  font-family: 'Courier New', Courier, monospace;
}

.video-placeholder {
  width: 300px;
  height: 225px;
  background-color: #111;
  border-radius: 8px;
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 1px solid #444;
}

.video-placeholder.camera-off {
  background-color: #2c2c2c;
}

.call-controls {
  margin-top: 2rem;
  display: flex;
  gap: 1.5rem;
}

.control-btn {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: none;
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: background-color 0.2s;
}

.control-btn:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.control-btn.toggled-off {
  background-color: rgba(255, 255, 255, 0.1);
  color: #aaa;
}

.end-call-btn {
  background-color: #e63946;
}
.end-call-btn:hover {
  background-color: #c42d39;
}
</style>