import io
import json
import os
import sys
import zipfile

from aqt import mw, gui_hooks
from aqt.utils import showWarning, qconnect
from aqt.qt import *

_addon_dir = os.path.dirname(__file__)
_vendor_dir = os.path.join(_addon_dir, "vendor")

_grpc_client = None
_core = None


def _wheel_score(filename, py_tag, machine):
    """Return compatibility score for a wheel filename. Higher = better. -1 = incompatible."""
    if not filename.endswith(".whl"):
        return -1
    # Pure-Python wheel works everywhere
    if "none-any" in filename:
        return 5
    # Must match our Python version
    if py_tag not in filename:
        return -1
    if sys.platform == "darwin":
        if "macosx" not in filename:
            return -1
        if machine == "arm64":
            if "arm64" in filename:
                return 20
            if "universal2" in filename:
                return 15
        else:  # x86_64
            if "x86_64" in filename:
                return 20
            if "universal2" in filename:
                return 15
    elif sys.platform == "linux":
        if "manylinux" not in filename and "linux" not in filename:
            return -1
        if machine in ("aarch64", "arm64"):
            if "aarch64" in filename:
                return 20
        else:  # x86_64
            if "x86_64" in filename:
                return 20
    elif sys.platform == "win32":
        if "win" not in filename:
            return -1
        if machine in ("AMD64", "x86_64"):
            if "win_amd64" in filename:
                return 20
            if "win32" in filename:
                return 10
        else:
            if "win32" in filename:
                return 20
    return -1


def _download_and_extract(url):
    import urllib.request
    data = urllib.request.urlopen(url, timeout=120).read()
    os.makedirs(_vendor_dir, exist_ok=True)
    with zipfile.ZipFile(io.BytesIO(data)) as zf:
        zf.extractall(_vendor_dir)


def _install_package(package, py_tag, machine):
    import urllib.request
    api_url = "https://pypi.org/pypi/{}/json".format(package)
    data = json.loads(urllib.request.urlopen(api_url, timeout=30).read().decode())

    # Try latest version first, then walk back through recent releases
    versions = [data["info"]["version"]] + list(data["releases"].keys())[-20:]
    for version in versions:
        files = data["releases"].get(version, [])
        best = None
        best_score = -1
        for f in files:
            s = _wheel_score(f["filename"], py_tag, machine)
            if s > best_score:
                best_score = s
                best = f
        if best and best_score >= 0:
            _download_and_extract(best["url"])
            return
    raise RuntimeError("No compatible wheel found for {} (py={}, arch={})".format(
        package, py_tag, machine))


def _ensure_grpc():
    if os.path.isdir(_vendor_dir) and _vendor_dir not in sys.path:
        sys.path.insert(0, _vendor_dir)
    try:
        import grpc  # noqa: F401
        return None
    except ImportError:
        pass

    import platform as _plat
    import traceback as _tb
    py_tag = "cp{}{}".format(sys.version_info.major, sys.version_info.minor)
    machine = _plat.machine()  # arm64 or x86_64

    try:
        _install_package("typing_extensions", py_tag, machine)
        _install_package("grpcio", py_tag, machine)
        _install_package("protobuf", py_tag, machine)
    except Exception:
        return _tb.format_exc()

    if _vendor_dir not in sys.path:
        sys.path.insert(0, _vendor_dir)
    try:
        import grpc  # noqa: F401
        return None
    except ImportError:
        return _tb.format_exc()


def _get_core():
    global _grpc_client, _core
    if _core is None:
        from .grpc_client import GrpcClient
        from .core import Core
        _grpc_client = GrpcClient()
        _core = Core(_grpc_client)
    return _core


def show_sengen():
    err = _ensure_grpc()
    if err:
        showWarning("Sengen: grpcio install failed:\n\n" + err)
        return
    try:
        from .gui import SengenDialog
        dlg = SengenDialog(_get_core(), mw)
        dlg.exec()
    except Exception:
        import traceback
        showWarning("Sengen error:\n\n" + traceback.format_exc())


def setup():
    action = QAction("Sengen", mw)
    qconnect(action.triggered, show_sengen)
    mw.form.menuTools.addAction(action)


gui_hooks.main_window_did_init.append(setup)
