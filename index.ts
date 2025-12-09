import { cpp } from "@/lib/languages/cpp";

Object.entries(cpp).forEach(([k, v]) => {
    console.log(k, " - ", v());
});
