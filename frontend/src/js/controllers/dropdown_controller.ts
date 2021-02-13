import { Controller } from "stimulus";
import { useClickOutside } from "stimulus-use";

export default class extends Controller {
  static targets = ["toggleable"];

  toggleableTarget: Element;
  toggleableTargets: Element[];
  hasToggleableTarget: boolean;

  connect() {
    useClickOutside(this);
  }

  toggle() {
    this.toggleableTarget.classList.toggle("hidden");
  }

  clickOutside(event: any) {
    if (this.toggleableTarget.classList.contains("hidden")) {
      return;
    }

    this.toggleableTarget.classList.add("hidden");
  }
}
