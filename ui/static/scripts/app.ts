interface Position {
  x: number;
  y: number;
}

interface HTMLElementWithStyle extends HTMLElement {
  style: CSSStyleDeclaration;
}

document.addEventListener("DOMContentLoaded", (): void => {
  const noBtn = document.getElementById("noBtn") as HTMLButtonElement;
  const yesBtn = document.getElementById("yesBtn") as HTMLButtonElement;
  const response = document.getElementById("response") as HTMLDivElement;
  const container = document.querySelector(".container") as HTMLDivElement;
  const buttonsContainer = document.querySelector(
    ".buttons-container",
  ) as HTMLDivElement;

  let countdownInterval: number | null = null;
  let isFirstHover: boolean = true;
  const originalNoBtnParent: HTMLElement | null = noBtn.parentElement;
  const originalNoBtnNextSibling: Element | null = noBtn.nextElementSibling;

  // Function to get random position within the viewport
  const getRandomPosition = (): Position => {
    const viewportWidth: number = window.innerWidth - noBtn.offsetWidth;
    const viewportHeight: number = window.innerHeight - noBtn.offsetHeight;

    return {
      x: Math.max(0, Math.floor(Math.random() * viewportWidth)),
      y: Math.max(0, Math.floor(Math.random() * viewportHeight)),
    };
  };

  // Function to create placeholder text
  const createPlaceholder = (): void => {
    const placeholder: HTMLDivElement = document.createElement("div");
    placeholder.textContent = "Don't click No please 💔";
    placeholder.className = "no-placeholder";
    placeholder.style.fontSize = "1rem";
    placeholder.style.color = "#ff6b6b";
    placeholder.style.display = "flex";
    placeholder.style.justifyContent = "center";
    placeholder.style.alignItems = "center";
    placeholder.style.minHeight = "38px";
    placeholder.style.width = "120px";
    // Insert the placeholder where the No button was
    noBtn.insertAdjacentElement("beforebegin", placeholder);
  };

  // Function to move the No button
  const moveNoButton = (): void => {
    const newPos: Position = getRandomPosition();
    noBtn.style.position = "fixed";
    noBtn.style.left = `${newPos.x}px`;
    noBtn.style.top = `${newPos.y}px`;
    noBtn.style.transition = "all 0.3s ease";
  };

  // Add one-time hover event to container
  const triggerFirstMove = function (): void {
    if (isFirstHover) {
      createPlaceholder();
      document.body.appendChild(noBtn); // Move to body for absolute positioning
      moveNoButton();
      isFirstHover = false;

      // Remove this event listener after first trigger
      container.removeEventListener("mouseenter", triggerFirstMove);

      // Add mouseover event to the No button for subsequent moves
      noBtn.addEventListener("mouseover", moveNoButton);
    }
  };

  container.addEventListener("mouseenter", triggerFirstMove);

  // Function to stop countdown and cancel shutdown
  const stopCountdown = (): void => {
    if (countdownInterval) {
      clearInterval(countdownInterval);
      countdownInterval = null;
      fetch("/cancel-shutdown", {
        method: "POST",
      }).catch((err: Error) => console.log("Error canceling shutdown:", err));
    }
  };

  // Function to start countdown
  const startCountdown = (): void => {
    let secondsLeft: number = 30;
    const countdownEl: HTMLElement | null =
      document.getElementById("countdown");

    countdownInterval = window.setInterval(() => {
      secondsLeft--;
      if (countdownEl) {
        countdownEl.textContent = `Closing in ${secondsLeft} seconds`;
      }
      if (secondsLeft <= 0) {
        if (countdownInterval) clearInterval(countdownInterval);
        window.close();
      }
    }, 1000);
  };

  // Function to trigger manual shutdown
  const triggerManualShutdown = (): void => {
    fetch("/trigger-shutdown", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ immediate: true }),
    })
      .then(() => {
        window.close();
      })
      .catch((err: Error) =>
        console.log("Error triggering manual shutdown:", err),
      );
  };

  // Handle "Yes" button click
  yesBtn.addEventListener("click", (): void => {
    // Remove the No button if it's in the document body
    if (noBtn.parentElement === document.body) {
      document.body.removeChild(noBtn);
    }

    container.innerHTML = `
            <h1>Yay! I knew you'd say yes! 🎉</h1>
            <div class="gif-container">
                <img src="/static/img/happy.gif" alt="Happy" class="happy-gif">
            </div>
            <p style="font-size: 1.5rem; margin-top: 2rem;">You've made me the happiest person! ❤️</p>
        `;

    // Add controls outside the main container
    const controls: HTMLDivElement = document.createElement("div");
    controls.innerHTML = `
            <div id="countdown" class="countdown">Closing in 30 seconds</div>
            <div class="control-buttons">
                <button class="btn control-btn" onclick="window.stopCountdown()">Stop</button>
                <button class="btn control-btn shutdown-btn" onclick="window.triggerManualShutdown()">Close</button>
            </div>
        `;
    document.body.appendChild(controls);

    // Make functions available globally
    (window as any).stopCountdown = stopCountdown;
    (window as any).triggerManualShutdown = triggerManualShutdown;

    // Start countdown and trigger shutdown
    startCountdown();
    fetch("/trigger-shutdown", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ immediate: false }),
    }).catch((err: Error) => console.log("Error triggering shutdown:", err));
  });
});
