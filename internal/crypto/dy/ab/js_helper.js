const fs = require('fs');
const path = require('path');

// Load the sign.js file
const signJsContent = fs.readFileSync(path.join(__dirname, 'sign.js'), 'utf8');

// Use eval to make the functions available.
eval(signJsContent);

// --- Monkey-patch for deterministic testing ---
const fixed_random_values = [0.123, 0.456, 0.789];
let random_idx = 0;
Math.random = () => {
    const val = fixed_random_values[random_idx % fixed_random_values.length];
    random_idx++;
    return val;
};

// A fixed date: 2023-03-15T12:00:00.000Z
const fixed_timestamp = 1678886400000; 
Date.now = () => fixed_timestamp;
// --- End of monkey-patch ---


// Test data (same as in the Python test)
const url_search_params = "device_platform=webapp&aid=6383&channel=channel_pc_web&update_version_code=170400&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1536&screen_height=864&browser_language=zh-CN&browser_platform=Win32&browser_name=Chrome&browser_version=123.0.0.0&browser_online=true&engine_name=Blink&engine_version=123.0.0.0&os_name=Windows&os_version=10&cpu_core_num=16&device_memory=8&platform=PC&downlink=10&effective_type=4g&round_trip_time=50&webid=7362810250930783783&msToken=VkDUvz1y24CppXSl80iFPr6ez-3FiizcwD7fI1OqBt6IICq9RWG7nCvxKb8IVi55mFd-wnqoNkXGnxHrikQb4PuKob5Q-YhDp5Um215JzlBszkUyiEvR";
const user_agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36";

// Generate the signature
const signature = sign_datail(url_search_params, user_agent);

// Print to stdout
console.log(signature);
