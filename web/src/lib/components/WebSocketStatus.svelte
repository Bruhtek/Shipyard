<script lang="ts">
    import WSDataStore from "$lib/websocket/MessageHandler.svelte.js";
    import { slide } from 'svelte/transition';
</script>

{#if WSDataStore.connectionStatus === 2}
    <div class="message green"
         in:slide={{duration: 500, axis: 'y'}}
         out:slide={{delay: 2000, duration: 500, axis: 'y'}}
    >
        Connection re-established!
    </div>
{/if}
{#if WSDataStore.connectionStatus === 2}
    <div class="message"
         in:slide={{duration: 500, axis: 'y'}}
         out:slide={{duration: 500, axis: 'y'}}
    >
        Disconnected from WebSocket. Attempting to connect...
    </div>
{/if}



<style>
    .message {
        position: absolute;
        top: 0;
        right: 0;
        width: 100%;
        background-color: var(--red-a0);
        color: var(--surface-tonal-a0);
        text-align: center;

        &.green {
            background-color: var(--green-a0);
        }
    }
</style>