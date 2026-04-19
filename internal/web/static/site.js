(() => {
  const activeClass = "is-open";
  const belowClass = "is-below";

  function positionTooltip(button, bubble) {
    if (!bubble) {
      return;
    }

    const viewportPadding = 16;

    bubble.classList.remove(belowClass);
    bubble.style.left = "50%";
    bubble.style.transform = "translateX(-50%)";

    const buttonRect = button.getBoundingClientRect();
    const bubbleRect = bubble.getBoundingClientRect();
    const maxLeft = window.innerWidth - viewportPadding - bubbleRect.width;
    const preferredLeft = buttonRect.left + buttonRect.width / 2 - bubbleRect.width / 2;
    const clampedLeft = Math.min(Math.max(preferredLeft, viewportPadding), maxLeft);

    bubble.style.left = `${clampedLeft - buttonRect.left}px`;
    bubble.style.transform = "none";

    if (bubble.getBoundingClientRect().top < 12) {
      bubble.classList.add(belowClass);
    }
  }

  async function ensureTooltip(button) {
    const bubble = button.querySelector(".tooltip-bubble");
    if (!bubble || bubble.dataset.loaded === "true") {
      return bubble;
    }

    const slug = button.dataset.tooltipSlug;
    bubble.hidden = false;
    bubble.innerHTML = '<div class="tooltip-card"><p>Loading...</p></div>';

    try {
      const response = await fetch(`/tooltips/${slug}`, {
        headers: { "X-Requested-With": "fetch" },
      });
      if (!response.ok) {
        throw new Error(`tooltip request failed: ${response.status}`);
      }
      bubble.innerHTML = await response.text();
      bubble.dataset.loaded = "true";
    } catch (_error) {
      bubble.innerHTML = '<div class="tooltip-card"><p>Tooltip unavailable.</p></div>';
    }

    return bubble;
  }

  function openTooltip(button) {
    button.classList.add(activeClass);
    ensureTooltip(button).then((bubble) => {
      if (bubble && button.classList.contains(activeClass)) {
        bubble.hidden = false;
        positionTooltip(button, bubble);
      }
    });
  }

  function closeTooltip(button) {
    button.classList.remove(activeClass);
    const bubble = button.querySelector(".tooltip-bubble");
    if (bubble) {
      bubble.hidden = true;
    }
  }

  function bindTooltip(button) {
    button.addEventListener("mouseenter", () => openTooltip(button));
    button.addEventListener("mouseleave", () => closeTooltip(button));
    button.addEventListener("focus", () => openTooltip(button));
    button.addEventListener("blur", () => closeTooltip(button));
    button.addEventListener("keydown", (event) => {
      if (event.key === "Escape") {
        closeTooltip(button);
        button.blur();
      }
    });
  }

  document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll(".tooltip-term").forEach(bindTooltip);
    window.addEventListener("resize", () => {
      document.querySelectorAll(`.tooltip-term.${activeClass}`).forEach((button) => {
        positionTooltip(button, button.querySelector(".tooltip-bubble"));
      });
    });
  });
})();
