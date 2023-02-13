<meta charset="UTF-8">
<title>提交结果</title>
<form method="POST">
    <p>请输入 uid <input type="text" name="uid" style="width: 10em;" value="<?php echo htmlspecialchars($_POST['uid']); ?>"></p>
    <p>请填写所有找到的 flag（使用空格分隔）: <input type="text" name="code" style="width: 40em;" value="<?php echo htmlspecialchars($_POST['code']); ?>"></p>
    <p><input type="submit"></p>
</form><?php

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $time = time();
    echo '<p>题目 A：', var_export(verify_key_2023_ganlv_flag((string)$_POST['code'], (string)$_POST['uid'], $time, 'flagA', '_2023challenge_52pojie_cn_', ['flag1{52pojiehappynewyear}', 'flag2{878a48f2}', 'flag3{06f95afd}', 'flag4{9cb91117}'], 2), true), '</p>';
    echo '<p>题目 B：', var_export(verify_key_2023_ganlv_flag((string)$_POST['code'], (string)$_POST['uid'], $time, 'flagB', '_happy_new_year_', ['flag5{eait}', 'flag6{590124}', 'flag7{5d06be63}', 'flag8{c394d7}'], 2), true), '</p>';
    echo '<p>题目 C：', var_export(verify_key_2023_ganlv_flag((string)$_POST['code'], (string)$_POST['uid'], $time, 'flagC', '_jwt_must_verify_sign_', ['flag9{21c5f8}', 'flag10{4a752b}', 'flag11{63418de7}', 'flag12{3ac97e24}'], 2), true), '</p>';
}

function verify_key_2023_ganlv_flag($flag, $uid, $timestamp, $dynamic_flag_name, $dynamic_secret, $static_flags, $need_static_flag_count) {
    $dynamic_flag = $dynamic_flag_name . '{' . substr(md5($uid . $dynamic_secret . floor($timestamp / 600)), 0, 8) . '}';
    $found_dynamic_flag = 0;
    if (false !== strpos($flag, $dynamic_flag)) {
        $found_dynamic_flag++;
    }
    $found_static_flag = 0;
    foreach ($static_flags as $static_flag) {
        if (false !== strpos($flag, $static_flag)) {
            $found_static_flag++;
        }
    }
    return [
        '动态 flag' => $dynamic_flag,
        '答对动态 flag 数量' => $found_dynamic_flag,
        '答对静态 flag 数量' => $found_static_flag,
    ];
}
